package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrPasswordResetTokenNotFound is returned by PasswordResetTokenRepository
// when no token matches the given hash.
var ErrPasswordResetTokenNotFound = errors.New("password reset token not found")

// PasswordResetToken is a persisted record of an issued password-reset
// token, keyed by a hash of the raw token (never the raw value itself, same
// rationale as RefreshToken) so the token can be looked up and consumed
// without the DB ever holding the value that was emailed to the user.
type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// IsValid reports whether the token is still usable — neither used nor
// expired.
func (t *PasswordResetToken) IsValid() bool {
	return t.UsedAt == nil && time.Now().Before(t.ExpiresAt)
}

// PasswordResetTokenRepository persists and retrieves PasswordResetToken
// records.
type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *PasswordResetToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	MarkUsed(ctx context.Context, tokenHash string) error
}
