package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

func newTestWorkspace() *domain.Workspace {
	return &domain.Workspace{
		Name: "Acme Inc.",
		Slug: "acme-" + uuid.NewString(),
	}
}

func TestWorkspaceRepository_CreateAndGetByID(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewWorkspaceRepository(tx)
		userRepo := NewUserRepository(tx)
		ctx := context.Background()

		owner := newTestUser()
		if err := userRepo.Create(ctx, owner); err != nil {
			t.Fatalf("creating owner user: %v", err)
		}

		workspace := newTestWorkspace()
		member := &domain.WorkspaceMember{UserID: owner.ID, Role: domain.RoleOwner}
		if err := repo.Create(ctx, workspace, member); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if workspace.ID == uuid.Nil {
			t.Fatal("Create() did not populate workspace ID")
		}
		if member.ID == uuid.Nil {
			t.Fatal("Create() did not populate owner membership ID")
		}
		if member.WorkspaceID != workspace.ID {
			t.Errorf("member.WorkspaceID = %v, want %v", member.WorkspaceID, workspace.ID)
		}

		got, err := repo.GetByID(ctx, workspace.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.Name != workspace.Name {
			t.Errorf("GetByID().Name = %q, want %q", got.Name, workspace.Name)
		}

		memberRepo := NewWorkspaceMemberRepository(tx)
		gotMember, err := memberRepo.GetByWorkspaceAndUser(ctx, workspace.ID, owner.ID)
		if err != nil {
			t.Fatalf("expected the owner membership row to exist, got error = %v", err)
		}
		if gotMember.Role != domain.RoleOwner {
			t.Errorf("Role = %v, want %v", gotMember.Role, domain.RoleOwner)
		}
	})
}

func TestWorkspaceRepository_GetByID_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewWorkspaceRepository(tx)

		_, err := repo.GetByID(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrWorkspaceNotFound) {
			t.Errorf("GetByID() error = %v, want %v", err, domain.ErrWorkspaceNotFound)
		}
	})
}

func TestWorkspaceRepository_Update(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewWorkspaceRepository(tx)
		userRepo := NewUserRepository(tx)
		ctx := context.Background()

		owner := newTestUser()
		if err := userRepo.Create(ctx, owner); err != nil {
			t.Fatalf("creating owner user: %v", err)
		}
		workspace := newTestWorkspace()
		if err := repo.Create(ctx, workspace, &domain.WorkspaceMember{UserID: owner.ID, Role: domain.RoleOwner}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		workspace.Name = "Acme Corp"
		workspace.LogoURL = "https://example.com/logo.png"
		if err := repo.Update(ctx, workspace); err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		got, err := repo.GetByID(ctx, workspace.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.Name != "Acme Corp" {
			t.Errorf("Name = %q, want %q", got.Name, "Acme Corp")
		}
		if got.LogoURL != "https://example.com/logo.png" {
			t.Errorf("LogoURL = %q, want %q", got.LogoURL, "https://example.com/logo.png")
		}
	})
}

func TestWorkspaceRepository_Update_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewWorkspaceRepository(tx)

		err := repo.Update(context.Background(), &domain.Workspace{ID: uuid.New(), Name: "x", Slug: "x-" + uuid.NewString()})
		if !errors.Is(err, domain.ErrWorkspaceNotFound) {
			t.Errorf("Update() error = %v, want %v", err, domain.ErrWorkspaceNotFound)
		}
	})
}

func TestWorkspaceRepository_SlugExists(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewWorkspaceRepository(tx)
		userRepo := NewUserRepository(tx)
		ctx := context.Background()

		owner := newTestUser()
		if err := userRepo.Create(ctx, owner); err != nil {
			t.Fatalf("creating owner user: %v", err)
		}
		workspace := newTestWorkspace()
		if err := repo.Create(ctx, workspace, &domain.WorkspaceMember{UserID: owner.ID, Role: domain.RoleOwner}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		exists, err := repo.SlugExists(ctx, workspace.Slug)
		if err != nil {
			t.Fatalf("SlugExists() error = %v", err)
		}
		if !exists {
			t.Error("SlugExists() = false, want true for a slug that was just created")
		}

		exists, err = repo.SlugExists(ctx, "no-such-slug-"+uuid.NewString())
		if err != nil {
			t.Fatalf("SlugExists() error = %v", err)
		}
		if exists {
			t.Error("SlugExists() = true, want false for a slug that doesn't exist")
		}
	})
}

func TestWorkspaceRepository_ListForUser(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewWorkspaceRepository(tx)
		userRepo := NewUserRepository(tx)
		ctx := context.Background()

		owner := newTestUser()
		if err := userRepo.Create(ctx, owner); err != nil {
			t.Fatalf("creating owner user: %v", err)
		}
		other := newTestUser()
		if err := userRepo.Create(ctx, other); err != nil {
			t.Fatalf("creating other user: %v", err)
		}

		workspace := newTestWorkspace()
		if err := repo.Create(ctx, workspace, &domain.WorkspaceMember{UserID: owner.ID, Role: domain.RoleOwner}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		otherWorkspace := newTestWorkspace()
		if err := repo.Create(ctx, otherWorkspace, &domain.WorkspaceMember{UserID: other.ID, Role: domain.RoleOwner}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.ListForUser(ctx, owner.ID)
		if err != nil {
			t.Fatalf("ListForUser() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 workspace, got %d", len(got))
		}
		if got[0].ID != workspace.ID {
			t.Errorf("ID = %v, want %v", got[0].ID, workspace.ID)
		}
	})
}
