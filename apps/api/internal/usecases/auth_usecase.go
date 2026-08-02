// Package usecases holds business logic. It depends only on domain (never
// on handlers or a specific repository implementation) per
// docs/adr/0002-clean-architecture-backend.md.
package usecases

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
	"github.com/foundryhq/foundryhq/apps/api/pkg/mailer"
)

// invalidCredentialsMessage is used for both "no such email" and "wrong
// password" so a login attempt can't be used to enumerate registered
// emails.
const invalidCredentialsMessage = "invalid email or password"

// invalidResetTokenMessage is used for both "no such token" and "expired or
// already-used token" so a reset attempt can't be used to distinguish the
// two.
const invalidResetTokenMessage = "invalid or expired reset token"

// invalidInviteTokenMessage mirrors invalidResetTokenMessage's
// enumeration-safety rationale for invite-activation tokens.
const invalidInviteTokenMessage = "invalid or expired invite token"

const minPasswordLength = 8

// maxPasswordLength bounds request size for the bcrypt-hashing call, not
// just usability — bcrypt silently truncates beyond 72 bytes anyway, but
// without a cap here the server would still spend a full bcrypt hash (a
// deliberately expensive operation) on an arbitrarily large attacker-supplied
// string first.
const maxPasswordLength = 72

// passwordResetTokenExpiry bounds how long an emailed reset link stays
// usable, short enough to limit exposure if the email is intercepted.
const passwordResetTokenExpiry = 1 * time.Hour

// inviteTokenExpiry bounds how long an emailed invite-activation link stays
// usable — longer than a password reset link since an invite is often
// acted on well after it's sent (vacations, backlog inboxes), not typically
// within the hour.
const inviteTokenExpiry = 7 * 24 * time.Hour

// AuthResult is returned by every AuthUsecase method that issues a session.
// RefreshToken/RefreshExpiresAt are exposed so the handler can set the
// httpOnly cookie (see adr/0004-jwt-access-refresh-tokens.md) — the usecase
// itself has no notion of cookies or HTTP.
type AuthResult struct {
	User             *domain.User
	AccessToken      string
	RefreshToken     string
	RefreshExpiresAt time.Time
}

// AuthUsecase implements registration, login, refresh, logout, password
// reset, and invite activation.
type AuthUsecase struct {
	userRepo               domain.UserRepository
	refreshTokenRepo       domain.RefreshTokenRepository
	passwordResetTokenRepo domain.PasswordResetTokenRepository
	inviteTokenRepo        domain.InviteTokenRepository
	workspaceMemberRepo    domain.WorkspaceMemberRepository
	jwtManager             *jwt.Manager
	mailer                 mailer.EmailSender
	appBaseURL             string
}

// NewAuthUsecase constructs an AuthUsecase. appBaseURL is the web app URL
// prefix used to build every emailed link (see config.AppBaseURL).
func NewAuthUsecase(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	passwordResetTokenRepo domain.PasswordResetTokenRepository,
	inviteTokenRepo domain.InviteTokenRepository,
	workspaceMemberRepo domain.WorkspaceMemberRepository,
	jwtManager *jwt.Manager,
	emailSender mailer.EmailSender,
	appBaseURL string,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:               userRepo,
		refreshTokenRepo:       refreshTokenRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		inviteTokenRepo:        inviteTokenRepo,
		workspaceMemberRepo:    workspaceMemberRepo,
		jwtManager:             jwtManager,
		mailer:                 emailSender,
		appBaseURL:             appBaseURL,
	}
}

// Register creates a new user and issues a session for them. An email that
// already belongs to ANY existing user — including a placeholder account
// created by WorkspaceUsecase.Invite, which has no password set yet — is
// rejected as already-registered, same as a fully-registered email. A
// placeholder can only be activated via AcceptInvite, which requires the
// emailed invite token as proof the caller actually owns that inbox;
// letting Register "claim" a placeholder by email alone let anyone who
// merely knew (or guessed) an invited address hijack that person's
// workspace membership before they ever signed up themselves.
func (u *AuthUsecase) Register(ctx context.Context, email, password string) (*AuthResult, error) {
	email = normalizeEmail(email)
	if email == "" {
		return nil, apperrors.Validation("email", "email is required")
	}
	if err := validatePasswordLength(password); err != nil {
		return nil, err
	}

	if _, err := u.userRepo.GetByEmail(ctx, email); err == nil {
		return nil, apperrors.Conflict("email already registered")
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, apperrors.Internal(fmt.Errorf("getting user: %w", err))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("hashing password: %w", err))
	}

	user := &domain.User{Email: email, PasswordHash: string(hash)}
	if err := u.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return nil, apperrors.Conflict("email already registered")
		}
		return nil, apperrors.Internal(fmt.Errorf("creating user: %w", err))
	}

	return u.issueSession(ctx, user)
}

// Login verifies email/password and issues a session on success.
func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	user, err := u.userRepo.GetByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, apperrors.Unauthorized(invalidCredentialsMessage)
		}
		return nil, apperrors.Internal(fmt.Errorf("getting user: %w", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apperrors.Unauthorized(invalidCredentialsMessage)
	}

	return u.issueSession(ctx, user)
}

// Refresh exchanges a valid, unrevoked refresh token for a new session,
// rotating the refresh token in the process (the presented one is revoked
// so it can't be replayed).
func (u *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	claims, err := u.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid refresh token")
	}

	stored, err := u.refreshTokenRepo.GetByTokenHash(ctx, hashToken(refreshToken))
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) {
			return nil, apperrors.Unauthorized("invalid refresh token")
		}
		return nil, apperrors.Internal(fmt.Errorf("getting refresh token: %w", err))
	}
	if !stored.IsValid() {
		return nil, apperrors.Unauthorized("refresh token expired or revoked")
	}

	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, apperrors.Unauthorized("invalid refresh token")
		}
		return nil, apperrors.Internal(fmt.Errorf("getting user: %w", err))
	}

	if err := u.refreshTokenRepo.Revoke(ctx, stored.TokenHash); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("revoking refresh token: %w", err))
	}

	return u.issueSession(ctx, user)
}

// Logout revokes the given refresh token. Revoking a token that's already
// gone isn't an error — the caller's session ends up logged out either way.
func (u *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	if err := u.refreshTokenRepo.Revoke(ctx, hashToken(refreshToken)); err != nil {
		return apperrors.Internal(fmt.Errorf("revoking refresh token: %w", err))
	}
	return nil
}

// ForgotPassword always succeeds, whether or not email is registered — same
// enumeration-safety principle as invalidCredentialsMessage. If the email
// matches a user, it generates a reset token, persists a hash of it, and
// emails the raw token as a link the user can follow to ResetPassword.
func (u *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	user, err := u.userRepo.GetByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return apperrors.Internal(fmt.Errorf("getting user: %w", err))
	}

	rawToken, err := generateSecureToken()
	if err != nil {
		return apperrors.Internal(fmt.Errorf("generating password reset token: %w", err))
	}

	if err := u.passwordResetTokenRepo.Create(ctx, &domain.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hashToken(rawToken),
		ExpiresAt: time.Now().Add(passwordResetTokenExpiry),
	}); err != nil {
		return apperrors.Internal(fmt.Errorf("persisting password reset token: %w", err))
	}

	resetLink := fmt.Sprintf("%s/auth/reset-password?token=%s", u.appBaseURL, rawToken)
	body := fmt.Sprintf(
		`<p>Click the link below to reset your FoundryHQ password. This link expires in %s.</p><p><a href="%s">%s</a></p>`,
		passwordResetTokenExpiry, resetLink, resetLink,
	)
	if err := u.mailer.Send(ctx, user.Email, "Reset your FoundryHQ password", body); err != nil {
		return apperrors.Internal(fmt.Errorf("sending password reset email: %w", err))
	}

	return nil
}

// ResetPassword validates token, sets newPassword as the account's new
// password, marks the token used so it can't be replayed, and revokes every
// existing refresh token for the account — forcing re-login everywhere on
// the assumption the old password may have been compromised.
func (u *AuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	if err := validatePasswordLength(newPassword); err != nil {
		return err
	}

	stored, err := u.passwordResetTokenRepo.GetByTokenHash(ctx, hashToken(token))
	if err != nil {
		if errors.Is(err, domain.ErrPasswordResetTokenNotFound) {
			return apperrors.Unauthorized(invalidResetTokenMessage)
		}
		return apperrors.Internal(fmt.Errorf("getting password reset token: %w", err))
	}
	if !stored.IsValid() {
		return apperrors.Unauthorized(invalidResetTokenMessage)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.Internal(fmt.Errorf("hashing password: %w", err))
	}

	if err := u.userRepo.UpdatePassword(ctx, stored.UserID, string(hash)); err != nil {
		return apperrors.Internal(fmt.Errorf("updating password: %w", err))
	}

	if err := u.passwordResetTokenRepo.MarkUsed(ctx, stored.TokenHash); err != nil {
		return apperrors.Internal(fmt.Errorf("marking password reset token used: %w", err))
	}

	if err := u.refreshTokenRepo.RevokeAllForUser(ctx, stored.UserID); err != nil {
		return apperrors.Internal(fmt.Errorf("revoking refresh tokens: %w", err))
	}

	return nil
}

// AcceptInvite activates the placeholder account behind an emailed invite
// token (see WorkspaceUsecase.Invite): sets the invitee's password, marks
// their workspace membership joined, marks the token used, and issues a
// session for them — the invite-activation equivalent of Register, but
// gated on proof of email ownership instead of the bare email address.
func (u *AuthUsecase) AcceptInvite(ctx context.Context, token, password string) (*AuthResult, error) {
	if err := validatePasswordLength(password); err != nil {
		return nil, err
	}

	stored, err := u.inviteTokenRepo.GetByTokenHash(ctx, hashToken(token))
	if err != nil {
		if errors.Is(err, domain.ErrInviteTokenNotFound) {
			return nil, apperrors.Unauthorized(invalidInviteTokenMessage)
		}
		return nil, apperrors.Internal(fmt.Errorf("getting invite token: %w", err))
	}
	if !stored.IsValid() {
		return nil, apperrors.Unauthorized(invalidInviteTokenMessage)
	}

	member, err := u.workspaceMemberRepo.GetByID(ctx, stored.WorkspaceMemberID)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("getting invited membership: %w", err))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("hashing password: %w", err))
	}
	if err := u.userRepo.UpdatePassword(ctx, member.UserID, string(hash)); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("activating invited account: %w", err))
	}
	if err := u.workspaceMemberRepo.MarkJoined(ctx, member.ID); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("marking membership joined: %w", err))
	}
	if err := u.inviteTokenRepo.MarkUsed(ctx, stored.TokenHash); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("marking invite token used: %w", err))
	}

	user, err := u.userRepo.GetByID(ctx, member.UserID)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("getting activated user: %w", err))
	}
	return u.issueSession(ctx, user)
}

// validatePasswordLength enforces both bounds on a caller-supplied password:
// short enough is a weak-password concern, long enough is a resource-cost
// concern (see maxPasswordLength).
func validatePasswordLength(password string) error {
	if len(password) < minPasswordLength {
		return apperrors.Validation("password", fmt.Sprintf("password must be at least %d characters", minPasswordLength))
	}
	if len(password) > maxPasswordLength {
		return apperrors.Validation("password", fmt.Sprintf("password must be at most %d characters", maxPasswordLength))
	}
	return nil
}

// issueSession generates a fresh access/refresh token pair for user and
// persists a hash of the refresh token so a later Logout/Refresh can find
// (and Logout can invalidate) it.
func (u *AuthUsecase) issueSession(ctx context.Context, user *domain.User) (*AuthResult, error) {
	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("generating access token: %w", err))
	}

	refreshToken, expiresAt, err := u.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("generating refresh token: %w", err))
	}

	if err := u.refreshTokenRepo.Create(ctx, &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: expiresAt,
	}); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("persisting refresh token: %w", err))
	}

	return &AuthResult{
		User:             user,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: expiresAt,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// hashToken returns a SHA-256 hex digest of token, which is what gets
// persisted (see domain.RefreshToken) — never the raw token, so a DB read
// alone can't be used to forge a session.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// generateSecureToken returns a random, URL-safe token — used for both the
// password-reset link (this file) and the invite-activation link
// (WorkspaceUsecase.Invite). Only its hash (see hashToken) is ever
// persisted; the raw value is emailed directly and never stored.
func generateSecureToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	return hex.EncodeToString(raw), nil
}
