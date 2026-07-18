// Package usecases holds business logic. It depends only on domain (never
// on handlers or a specific repository implementation) per
// docs/adr/0002-clean-architecture-backend.md.
package usecases

import (
	"context"
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
)

// invalidCredentialsMessage is used for both "no such email" and "wrong
// password" so a login attempt can't be used to enumerate registered
// emails.
const invalidCredentialsMessage = "invalid email or password"

const minPasswordLength = 8

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

// AuthUsecase implements registration, login, refresh, and logout.
type AuthUsecase struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	jwtManager       *jwt.Manager
}

// NewAuthUsecase constructs an AuthUsecase.
func NewAuthUsecase(userRepo domain.UserRepository, refreshTokenRepo domain.RefreshTokenRepository, jwtManager *jwt.Manager) *AuthUsecase {
	return &AuthUsecase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
	}
}

// Register creates a new user and issues a session for them.
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
