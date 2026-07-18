package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// mustCreateTestUser satisfies refresh_tokens' user_id foreign key.
func mustCreateTestUser(t *testing.T, tx *gorm.DB) *domain.User {
	t.Helper()
	user := newTestUser()
	if err := NewUserRepository(tx).Create(context.Background(), user); err != nil {
		t.Fatalf("creating test user: %v", err)
	}
	return user
}

func TestRefreshTokenRepository_CreateAndGetByTokenHash(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		user := mustCreateTestUser(t, tx)
		repo := NewRefreshTokenRepository(tx)
		ctx := context.Background()

		token := &domain.RefreshToken{
			UserID:    user.ID,
			TokenHash: "hash-" + uuid.NewString(),
			ExpiresAt: time.Now().Add(168 * time.Hour),
		}
		if err := repo.Create(ctx, token); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if token.ID == uuid.Nil {
			t.Fatal("Create() did not populate ID")
		}

		got, err := repo.GetByTokenHash(ctx, token.TokenHash)
		if err != nil {
			t.Fatalf("GetByTokenHash() error = %v", err)
		}
		if got.UserID != user.ID {
			t.Errorf("GetByTokenHash().UserID = %v, want %v", got.UserID, user.ID)
		}
		if !got.IsValid() {
			t.Error("IsValid() = false, want true for a freshly created token")
		}
	})
}

func TestRefreshTokenRepository_GetByTokenHash_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewRefreshTokenRepository(tx)

		_, err := repo.GetByTokenHash(context.Background(), "nonexistent-hash")
		if !errors.Is(err, domain.ErrRefreshTokenNotFound) {
			t.Errorf("GetByTokenHash() error = %v, want %v", err, domain.ErrRefreshTokenNotFound)
		}
	})
}

func TestRefreshTokenRepository_Revoke(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		user := mustCreateTestUser(t, tx)
		repo := NewRefreshTokenRepository(tx)
		ctx := context.Background()

		token := &domain.RefreshToken{
			UserID:    user.ID,
			TokenHash: "hash-" + uuid.NewString(),
			ExpiresAt: time.Now().Add(168 * time.Hour),
		}
		if err := repo.Create(ctx, token); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if err := repo.Revoke(ctx, token.TokenHash); err != nil {
			t.Fatalf("Revoke() error = %v", err)
		}

		got, err := repo.GetByTokenHash(ctx, token.TokenHash)
		if err != nil {
			t.Fatalf("GetByTokenHash() error = %v", err)
		}
		if got.IsValid() {
			t.Error("IsValid() = true, want false after Revoke")
		}
	})
}

func TestRefreshTokenRepository_Revoke_NonexistentIsNoop(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewRefreshTokenRepository(tx)

		if err := repo.Revoke(context.Background(), "never-existed"); err != nil {
			t.Errorf("Revoke() on nonexistent token error = %v, want nil", err)
		}
	})
}
