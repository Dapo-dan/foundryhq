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

const minPasswordLength = 8

// passwordResetTokenExpiry bounds how long an emailed reset link stays
// usable, short enough to limit exposure if the email is intercepted.
const passwordResetTokenExpiry = 1 * time.Hour

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

// AuthUsecase implements registration, login, refresh, logout, and password
// reset.
type AuthUsecase struct {
	userRepo               domain.UserRepository
	refreshTokenRepo       domain.RefreshTokenRepository
	passwordResetTokenRepo domain.PasswordResetTokenRepository
	jwtManager             *jwt.Manager
	mailer                 mailer.EmailSender
	passwordResetURLBase   string
}

// NewAuthUsecase constructs an AuthUsecase. passwordResetURLBase is the web
// app URL prefix used to build the emailed reset link (see
// config.PasswordResetURLBase).
func NewAuthUsecase(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	passwordResetTokenRepo domain.PasswordResetTokenRepository,
	jwtManager *jwt.Manager,
	emailSender mailer.EmailSender,
	passwordResetURLBase string,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:               userRepo,
		refreshTokenRepo:       refreshTokenRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		jwtManager:             jwtManager,
		mailer:                 emailSender,
		passwordResetURLBase:   passwordResetURLBase,
	}
}

// Register creates a new user and issues a session for them. If email
// already belongs to a placeholder account (invited into a workspace via
// WorkspaceUsecase.Invite before ever signing up — see that method's doc
// comment) it claims that account by setting its password instead of
// rejecting the registration, since a placeholder has no password for the
// caller to have "already registered" with.
func (u *AuthUsecase) Register(ctx context.Context, email, password string) (*AuthResult, error) {
	email = normalizeEmail(email)
	if email == "" {
		return nil, apperrors.Validation("email", "email is required")
	}
	if len(password) < minPasswordLength {
		return nil, apperrors.Validation("password", fmt.Sprintf("password must be at least %d characters", minPasswordLength))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("hashing password: %w", err))
	}

	existing, err := u.userRepo.GetByEmail(ctx, email)
	if err == nil {
		if existing.PasswordHash != "" {
			return nil, apperrors.Conflict("email already registered")
		}
		if err := u.userRepo.UpdatePassword(ctx, existing.ID, string(hash)); err != nil {
			return nil, apperrors.Internal(fmt.Errorf("claiming placeholder account: %w", err))
		}
		existing.PasswordHash = string(hash)
		return u.issueSession(ctx, existing)
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, apperrors.Internal(fmt.Errorf("getting user: %w", err))
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

	rawToken, err := generateResetToken()
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

	resetLink := fmt.Sprintf("%s/auth/reset-password?token=%s", u.passwordResetURLBase, rawToken)
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
	if len(newPassword) < minPasswordLength {
		return apperrors.Validation("password", fmt.Sprintf("password must be at least %d characters", minPasswordLength))
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

// generateResetToken returns a random, URL-safe password-reset token. Only
// its hash (see hashToken) is ever persisted; the raw value is emailed
// directly to the user and never stored.
func generateResetToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	return hex.EncodeToString(raw), nil
}
