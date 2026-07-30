package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

func TestWorkspaceMemberRepository_CreateAndGetByWorkspaceAndUser(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		userRepo := NewUserRepository(tx)
		workspaceRepo := NewWorkspaceRepository(tx)
		memberRepo := NewWorkspaceMemberRepository(tx)
		ctx := context.Background()

		owner := newTestUser()
		if err := userRepo.Create(ctx, owner); err != nil {
			t.Fatalf("creating owner user: %v", err)
		}
		workspace := newTestWorkspace()
		if err := workspaceRepo.Create(ctx, workspace, &domain.WorkspaceMember{UserID: owner.ID, Role: domain.RoleOwner}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		invitee := newTestUser()
		if err := userRepo.Create(ctx, invitee); err != nil {
			t.Fatalf("creating invitee user: %v", err)
		}
		member := &domain.WorkspaceMember{WorkspaceID: workspace.ID, UserID: invitee.ID, Role: domain.RoleMember}
		if err := memberRepo.Create(ctx, member); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if member.ID == uuid.Nil {
			t.Fatal("Create() did not populate ID")
		}
		if member.JoinedAt != nil {
			t.Error("a freshly invited member should have a nil JoinedAt")
		}

		got, err := memberRepo.GetByWorkspaceAndUser(ctx, workspace.ID, invitee.ID)
		if err != nil {
			t.Fatalf("GetByWorkspaceAndUser() error = %v", err)
		}
		if got.Role != domain.RoleMember {
			t.Errorf("Role = %v, want %v", got.Role, domain.RoleMember)
		}
	})
}

func TestWorkspaceMemberRepository_GetByWorkspaceAndUser_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewWorkspaceMemberRepository(tx)

		_, err := repo.GetByWorkspaceAndUser(context.Background(), uuid.New(), uuid.New())
		if !errors.Is(err, domain.ErrWorkspaceMemberNotFound) {
			t.Errorf("GetByWorkspaceAndUser() error = %v, want %v", err, domain.ErrWorkspaceMemberNotFound)
		}
	})
}

func TestWorkspaceMemberRepository_ListByWorkspaceIDWithUser(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		userRepo := NewUserRepository(tx)
		workspaceRepo := NewWorkspaceRepository(tx)
		memberRepo := NewWorkspaceMemberRepository(tx)
		ctx := context.Background()

		owner := newTestUser()
		if err := userRepo.Create(ctx, owner); err != nil {
			t.Fatalf("creating owner user: %v", err)
		}
		workspace := newTestWorkspace()
		if err := workspaceRepo.Create(ctx, workspace, &domain.WorkspaceMember{UserID: owner.ID, Role: domain.RoleOwner}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		invitee := newTestUser()
		if err := userRepo.Create(ctx, invitee); err != nil {
			t.Fatalf("creating invitee user: %v", err)
		}
		if err := memberRepo.Create(ctx, &domain.WorkspaceMember{WorkspaceID: workspace.ID, UserID: invitee.ID, Role: domain.RoleMember}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		members, err := memberRepo.ListByWorkspaceIDWithUser(ctx, workspace.ID)
		if err != nil {
			t.Fatalf("ListByWorkspaceIDWithUser() error = %v", err)
		}
		if len(members) != 2 {
			t.Fatalf("expected 2 members, got %d", len(members))
		}

		emails := map[string]bool{}
		for _, m := range members {
			emails[m.Email] = true
		}
		if !emails[owner.Email] || !emails[invitee.Email] {
			t.Errorf("expected both member emails present, got %v", emails)
		}
	})
}

func TestWorkspaceMemberRepository_UpdateRole(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		userRepo := NewUserRepository(tx)
		workspaceRepo := NewWorkspaceRepository(tx)
		memberRepo := NewWorkspaceMemberRepository(tx)
		ctx := context.Background()

		owner := newTestUser()
		if err := userRepo.Create(ctx, owner); err != nil {
			t.Fatalf("creating owner user: %v", err)
		}
		workspace := newTestWorkspace()
		ownerMember := &domain.WorkspaceMember{UserID: owner.ID, Role: domain.RoleOwner}
		if err := workspaceRepo.Create(ctx, workspace, ownerMember); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if err := memberRepo.UpdateRole(ctx, ownerMember.ID, domain.RoleMember); err != nil {
			t.Fatalf("UpdateRole() error = %v", err)
		}

		got, err := memberRepo.GetByWorkspaceAndUser(ctx, workspace.ID, owner.ID)
		if err != nil {
			t.Fatalf("GetByWorkspaceAndUser() error = %v", err)
		}
		if got.Role != domain.RoleMember {
			t.Errorf("Role = %v, want %v", got.Role, domain.RoleMember)
		}
	})
}

func TestWorkspaceMemberRepository_UpdateRole_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewWorkspaceMemberRepository(tx)

		err := repo.UpdateRole(context.Background(), uuid.New(), domain.RoleMember)
		if !errors.Is(err, domain.ErrWorkspaceMemberNotFound) {
			t.Errorf("UpdateRole() error = %v, want %v", err, domain.ErrWorkspaceMemberNotFound)
		}
	})
}
