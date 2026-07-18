package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// refreshTokenModel is the GORM row shape for the refresh_tokens table (see
// migration 000002_refresh_tokens).
type refreshTokenModel struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"column:user_id"`
	TokenHash string     `gorm:"column:token_hash"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func (refreshTokenModel) TableName() string { return "refresh_tokens" }

func (m refreshTokenModel) toDomain() *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		RevokedAt: m.RevokedAt,
		CreatedAt: m.CreatedAt,
	}
}

// RefreshTokenRepository implements domain.RefreshTokenRepository on top of
// GORM/Postgres.
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository constructs a RefreshTokenRepository.
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create inserts token, generating its ID and CreatedAt via the table's
// column defaults.
func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	model := &refreshTokenModel{
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("creating refresh token: %w", err)
	}

	token.ID = model.ID
	token.CreatedAt = model.CreatedAt
	return nil
}

// GetByTokenHash returns the token record matching tokenHash, or
// domain.ErrRefreshTokenNotFound. It does not filter by validity — callers
// should check RefreshToken.IsValid(), since what counts as valid is a
// usecase-level decision, not a query concern.
func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	var model refreshTokenModel
	if err := r.db.WithContext(ctx).First(&model, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("getting refresh token: %w", err)
	}
	return model.toDomain(), nil
}

// Revoke marks the token matching tokenHash as revoked. It's a no-op (not
// an error) if no such token exists — logging out a session that's already
// gone shouldn't fail the request.
func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	err := r.db.WithContext(ctx).
		Model(&refreshTokenModel{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", time.Now()).Error
	if err != nil {
		return fmt.Errorf("revoking refresh token: %w", err)
	}
	return nil
}
