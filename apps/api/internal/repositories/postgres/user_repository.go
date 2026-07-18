// Package postgres implements the domain repository interfaces on top of
// GORM. Models here are separate from domain entities (see
// docs/adr/0003-repository-pattern.md) so the domain package stays free of
// GORM struct tags and infra concerns.
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

// userModel is the GORM row shape for the users table (see migration
// 000001_init_schema). oauth_provider/oauth_id exist in the schema for the
// v1.1+ OAuth flow (docs/mvp.md) but aren't mapped here — nothing in the
// current codebase reads or writes them yet.
type userModel struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (userModel) TableName() string { return "users" }

func (m userModel) toDomain() *domain.User {
	return &domain.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func userModelFromDomain(u *domain.User) *userModel {
	return &userModel{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
}

// UserRepository implements domain.UserRepository on top of GORM/Postgres.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository constructs a UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts user, generating its ID and timestamps via the table's
// column defaults. Returns domain.ErrEmailAlreadyExists if the email is
// already registered.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := userModelFromDomain(user)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("creating user: %w", err)
	}

	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt
	return nil
}

// GetByID returns the user with the given ID, or domain.ErrUserNotFound.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("getting user %s: %w", id, err)
	}
	return model.toDomain(), nil
}

// GetByEmail returns the user with the given email, or domain.ErrUserNotFound.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).First(&model, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("getting user by email: %w", err)
	}
	return model.toDomain(), nil
}
