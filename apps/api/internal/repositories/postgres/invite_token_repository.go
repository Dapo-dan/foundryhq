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

// inviteTokenModel is the GORM row shape for the invite_tokens table (see
// migration 000007_invite_tokens).
type inviteTokenModel struct {
	ID                uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkspaceMemberID uuid.UUID  `gorm:"column:workspace_member_id"`
	TokenHash         string     `gorm:"column:token_hash"`
	ExpiresAt         time.Time  `gorm:"column:expires_at"`
	UsedAt            *time.Time `gorm:"column:used_at"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
}

func (inviteTokenModel) TableName() string { return "invite_tokens" }

func (m inviteTokenModel) toDomain() *domain.InviteToken {
	return &domain.InviteToken{
		ID:                m.ID,
		WorkspaceMemberID: m.WorkspaceMemberID,
		TokenHash:         m.TokenHash,
		ExpiresAt:         m.ExpiresAt,
		UsedAt:            m.UsedAt,
		CreatedAt:         m.CreatedAt,
	}
}

// InviteTokenRepository implements domain.InviteTokenRepository on top of
// GORM/Postgres.
type InviteTokenRepository struct {
	db *gorm.DB
}

// NewInviteTokenRepository constructs an InviteTokenRepository.
func NewInviteTokenRepository(db *gorm.DB) *InviteTokenRepository {
	return &InviteTokenRepository{db: db}
}

// Create inserts token, generating its ID and CreatedAt via the table's
// column defaults.
func (r *InviteTokenRepository) Create(ctx context.Context, token *domain.InviteToken) error {
	model := &inviteTokenModel{
		WorkspaceMemberID: token.WorkspaceMemberID,
		TokenHash:         token.TokenHash,
		ExpiresAt:         token.ExpiresAt,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("creating invite token: %w", err)
	}

	token.ID = model.ID
	token.CreatedAt = model.CreatedAt
	return nil
}

// GetByTokenHash returns the token record matching tokenHash, or
// domain.ErrInviteTokenNotFound. It does not filter by validity — callers
// should check InviteToken.IsValid().
func (r *InviteTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.InviteToken, error) {
	var model inviteTokenModel
	if err := r.db.WithContext(ctx).First(&model, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInviteTokenNotFound
		}
		return nil, fmt.Errorf("getting invite token: %w", err)
	}
	return model.toDomain(), nil
}

// MarkUsed marks the token matching tokenHash as used. It's a no-op (not an
// error) if no such token exists.
func (r *InviteTokenRepository) MarkUsed(ctx context.Context, tokenHash string) error {
	err := r.db.WithContext(ctx).
		Model(&inviteTokenModel{}).
		Where("token_hash = ?", tokenHash).
		Update("used_at", time.Now()).Error
	if err != nil {
		return fmt.Errorf("marking invite token used: %w", err)
	}
	return nil
}
