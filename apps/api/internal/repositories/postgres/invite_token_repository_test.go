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

// mustCreatePlaceholderMember seeds a workspace (with its own owner) plus a
// placeholder member row — the shape WorkspaceUsecase.Invite creates for an
// email that hasn't signed up yet — for invite-token tests to hang a token
// off of.
func mustCreatePlaceholderMember(t *testing.T, tx *gorm.DB) *domain.WorkspaceMember {
	t.Helper()
	workspace := newTestProjectWorkspace(t, tx)

	placeholder := &domain.User{Email: "invited-" + uuid.NewString() + "@example.com"}
	if err := NewUserRepository(tx).Create(context.Background(), placeholder); err != nil {
		t.Fatalf("creating placeholder user: %v", err)
	}

	member := &domain.WorkspaceMember{WorkspaceID: workspace.ID, UserID: placeholder.ID, Role: domain.RoleMember}
	if err := NewWorkspaceMemberRepository(tx).Create(context.Background(), member); err != nil {
		t.Fatalf("creating placeholder membership: %v", err)
	}
	return member
}

func TestInviteTokenRepository_CreateAndGetByTokenHash(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		member := mustCreatePlaceholderMember(t, tx)
		repo := NewInviteTokenRepository(tx)
		ctx := context.Background()

		token := &domain.InviteToken{
			WorkspaceMemberID: member.ID,
			TokenHash:         "hash-" + uuid.NewString(),
			ExpiresAt:         time.Now().Add(7 * 24 * time.Hour),
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
		if got.WorkspaceMemberID != member.ID {
			t.Errorf("GetByTokenHash().WorkspaceMemberID = %v, want %v", got.WorkspaceMemberID, member.ID)
		}
		if !got.IsValid() {
			t.Error("IsValid() = false, want true for a freshly created token")
		}
	})
}

func TestInviteTokenRepository_GetByTokenHash_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewInviteTokenRepository(tx)

		_, err := repo.GetByTokenHash(context.Background(), "nonexistent-hash")
		if !errors.Is(err, domain.ErrInviteTokenNotFound) {
			t.Errorf("GetByTokenHash() error = %v, want %v", err, domain.ErrInviteTokenNotFound)
		}
	})
}

func TestInviteTokenRepository_MarkUsed(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		member := mustCreatePlaceholderMember(t, tx)
		repo := NewInviteTokenRepository(tx)
		ctx := context.Background()

		token := &domain.InviteToken{
			WorkspaceMemberID: member.ID,
			TokenHash:         "hash-" + uuid.NewString(),
			ExpiresAt:         time.Now().Add(7 * 24 * time.Hour),
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

func TestInviteTokenRepository_MarkUsed_NonexistentIsNoop(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewInviteTokenRepository(tx)

		if err := repo.MarkUsed(context.Background(), "never-existed"); err != nil {
			t.Errorf("MarkUsed() on nonexistent token error = %v, want nil", err)
		}
	})
}
