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

func TestPasswordResetTokenRepository_CreateAndGetByTokenHash(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		user := mustCreateTestUser(t, tx)
		repo := NewPasswordResetTokenRepository(tx)
		ctx := context.Background()

		token := &domain.PasswordResetToken{
			UserID:    user.ID,
			TokenHash: "hash-" + uuid.NewString(),
			ExpiresAt: time.Now().Add(time.Hour),
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

func TestPasswordResetTokenRepository_GetByTokenHash_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewPasswordResetTokenRepository(tx)

		_, err := repo.GetByTokenHash(context.Background(), "nonexistent-hash")
		if !errors.Is(err, domain.ErrPasswordResetTokenNotFound) {
			t.Errorf("GetByTokenHash() error = %v, want %v", err, domain.ErrPasswordResetTokenNotFound)
		}
	})
}

func TestPasswordResetTokenRepository_MarkUsed(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		user := mustCreateTestUser(t, tx)
		repo := NewPasswordResetTokenRepository(tx)
		ctx := context.Background()

		token := &domain.PasswordResetToken{
			UserID:    user.ID,
			TokenHash: "hash-" + uuid.NewString(),
			ExpiresAt: time.Now().Add(time.Hour),
		}
		if err := repo.Create(ctx, token); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if err := repo.MarkUsed(ctx, token.TokenHash); err != nil {
			t.Fatalf("MarkUsed() error = %v", err)
		}

		got, err := repo.GetByTokenHash(ctx, token.TokenHash)
		if err != nil {
			t.Fatalf("GetByTokenHash() error = %v", err)
		}
		if got.IsValid() {
			t.Error("IsValid() = true, want false after MarkUsed")
		}
	})
}

func TestPasswordResetTokenRepository_MarkUsed_NonexistentIsNoop(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewPasswordResetTokenRepository(tx)

		if err := repo.MarkUsed(context.Background(), "never-existed"); err != nil {
			t.Errorf("MarkUsed() on nonexistent token error = %v, want nil", err)
		}
	})
}
