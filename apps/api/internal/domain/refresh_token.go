package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrRefreshTokenNotFound is returned by RefreshTokenRepository when no
// token matches the given hash.
var ErrRefreshTokenNotFound = errors.New("refresh token not found")

// RefreshToken is a persisted record of an issued refresh token, keyed by a
// hash of the raw token (never the raw value itself, see
// adr/0004-jwt-access-refresh-tokens.md) so logout can explicitly revoke it.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// IsValid reports whether the token is still usable — neither revoked nor
// expired. What counts as "valid" is a business rule, so it lives here
// rather than as a query filter in the repository.
func (t *RefreshToken) IsValid() bool {
	return t.RevokedAt == nil && time.Now().Before(t.ExpiresAt)
}

// RefreshTokenRepository persists and retrieves RefreshToken records.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
}
