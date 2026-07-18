package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

func newTestUser() *domain.User {
	return &domain.User{
		Email:        "user-" + uuid.NewString() + "@example.com",
		PasswordHash: "hashed-password",
	}
}

func TestUserRepository_CreateAndGetByID(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewUserRepository(tx)
		ctx := context.Background()

		user := newTestUser()
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if user.ID == uuid.Nil {
			t.Fatal("Create() did not populate ID")
		}

		got, err := repo.GetByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.Email != user.Email {
			t.Errorf("GetByID().Email = %q, want %q", got.Email, user.Email)
		}
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewUserRepository(tx)
		ctx := context.Background()

		user := newTestUser()
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.GetByEmail(ctx, user.Email)
		if err != nil {
			t.Fatalf("GetByEmail() error = %v", err)
		}
		if got.ID != user.ID {
			t.Errorf("GetByEmail().ID = %v, want %v", got.ID, user.ID)
		}
	})
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewUserRepository(tx)
		ctx := context.Background()

		user := newTestUser()
		if err := repo.Create(ctx, user); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		duplicate := &domain.User{Email: user.Email, PasswordHash: "another-hash"}
		err := repo.Create(ctx, duplicate)
		if !errors.Is(err, domain.ErrEmailAlreadyExists) {
			t.Errorf("Create() error = %v, want %v", err, domain.ErrEmailAlreadyExists)
		}
	})
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewUserRepository(tx)

		_, err := repo.GetByID(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("GetByID() error = %v, want %v", err, domain.ErrUserNotFound)
		}
	})
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewUserRepository(tx)

		_, err := repo.GetByEmail(context.Background(), "nobody@example.com")
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("GetByEmail() error = %v, want %v", err, domain.ErrUserNotFound)
		}
	})
}
