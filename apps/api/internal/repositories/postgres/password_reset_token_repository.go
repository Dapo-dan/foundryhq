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

// passwordResetTokenModel is the GORM row shape for the
// password_reset_tokens table (see migration
// 000004_password_reset_tokens).
type passwordResetTokenModel struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"column:user_id"`
	TokenHash string     `gorm:"column:token_hash"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func (passwordResetTokenModel) TableName() string { return "password_reset_tokens" }

func (m passwordResetTokenModel) toDomain() *domain.PasswordResetToken {
	return &domain.PasswordResetToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
		CreatedAt: m.CreatedAt,
	}
}

// PasswordResetTokenRepository implements domain.PasswordResetTokenRepository
// on top of GORM/Postgres.
type PasswordResetTokenRepository struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository constructs a PasswordResetTokenRepository.
func NewPasswordResetTokenRepository(db *gorm.DB) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{db: db}
}

// Create inserts token, generating its ID and CreatedAt via the table's
// column defaults.
func (r *PasswordResetTokenRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	model := &passwordResetTokenModel{
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("creating password reset token: %w", err)
	}

	token.ID = model.ID
	token.CreatedAt = model.CreatedAt
	return nil
}

// GetByTokenHash returns the token record matching tokenHash, or
// domain.ErrPasswordResetTokenNotFound. It does not filter by validity —
// callers should check PasswordResetToken.IsValid().
func (r *PasswordResetTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordResetToken, error) {
	var model passwordResetTokenModel
	if err := r.db.WithContext(ctx).First(&model, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPasswordResetTokenNotFound
		}
		return nil, fmt.Errorf("getting password reset token: %w", err)
	}
	return model.toDomain(), nil
}

// MarkUsed marks the token matching tokenHash as used. It's a no-op (not an
// error) if no such token exists.
func (r *PasswordResetTokenRepository) MarkUsed(ctx context.Context, tokenHash string) error {
	err := r.db.WithContext(ctx).
		Model(&passwordResetTokenModel{}).
		Where("token_hash = ?", tokenHash).
		Update("used_at", time.Now()).Error
	if err != nil {
		return fmt.Errorf("marking password reset token used: %w", err)
	}
	return nil
}
