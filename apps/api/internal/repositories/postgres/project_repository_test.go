package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// newTestProjectWorkspace creates and returns a workspace (with an owner
// membership) that tests can hang projects off of.
func newTestProjectWorkspace(t *testing.T, tx *gorm.DB) *domain.Workspace {
	t.Helper()
	userRepo := NewUserRepository(tx)
	workspaceRepo := NewWorkspaceRepository(tx)
	ctx := context.Background()

	owner := newTestUser()
	if err := userRepo.Create(ctx, owner); err != nil {
		t.Fatalf("creating owner user: %v", err)
	}
	workspace := newTestWorkspace()
	if err := workspaceRepo.Create(ctx, workspace, &domain.WorkspaceMember{UserID: owner.ID, Role: domain.RoleOwner}); err != nil {
		t.Fatalf("creating workspace: %v", err)
	}
	return workspace
}

func TestProjectRepository_CreateAndGetByID(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace := newTestProjectWorkspace(t, tx)
		repo := NewProjectRepository(tx)
		ctx := context.Background()

		description := "Q3 roadmap"
		project := &domain.Project{WorkspaceID: workspace.ID, Name: "Launch", Description: &description}
		if err := repo.Create(ctx, project); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if project.ID == uuid.Nil {
			t.Fatal("Create() did not populate ID")
		}

		got, err := repo.GetByID(ctx, project.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.Name != "Launch" {
			t.Errorf("Name = %q, want %q", got.Name, "Launch")
		}
		if got.Description == nil || *got.Description != description {
			t.Errorf("Description = %v, want %q", got.Description, description)
		}
	})
}

func TestProjectRepository_GetByID_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewProjectRepository(tx)

		_, err := repo.GetByID(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrProjectNotFound) {
			t.Errorf("GetByID() error = %v, want %v", err, domain.ErrProjectNotFound)
		}
	})
}

func TestProjectRepository_ListByWorkspaceID(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace := newTestProjectWorkspace(t, tx)
		otherWorkspace := newTestProjectWorkspace(t, tx)
		repo := NewProjectRepository(tx)
		ctx := context.Background()

		if err := repo.Create(ctx, &domain.Project{WorkspaceID: workspace.ID, Name: "Launch"}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if err := repo.Create(ctx, &domain.Project{WorkspaceID: otherWorkspace.ID, Name: "Other"}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.ListByWorkspaceID(ctx, workspace.ID)
		if err != nil {
			t.Fatalf("ListByWorkspaceID() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 project, got %d", len(got))
		}
		if got[0].Name != "Launch" {
			t.Errorf("Name = %q, want %q", got[0].Name, "Launch")
		}
	})
}

func TestProjectRepository_Update(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace := newTestProjectWorkspace(t, tx)
		repo := NewProjectRepository(tx)
		ctx := context.Background()

		project := &domain.Project{WorkspaceID: workspace.ID, Name: "Launch"}
		if err := repo.Create(ctx, project); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		description := "Now with a description"
		project.Name = "Launch v2"
		project.Description = &description
		if err := repo.Update(ctx, project); err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		got, err := repo.GetByID(ctx, project.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.Name != "Launch v2" {
			t.Errorf("Name = %q, want %q", got.Name, "Launch v2")
		}
		if got.Description == nil || *got.Description != description {
			t.Errorf("Description = %v, want %q", got.Description, description)
		}
	})
}

func TestProjectRepository_Update_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewProjectRepository(tx)

		err := repo.Update(context.Background(), &domain.Project{ID: uuid.New(), Name: "x"})
		if !errors.Is(err, domain.ErrProjectNotFound) {
			t.Errorf("Update() error = %v, want %v", err, domain.ErrProjectNotFound)
		}
	})
}

func TestProjectRepository_Delete_ExcludesFromFutureQueries(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace := newTestProjectWorkspace(t, tx)
		repo := NewProjectRepository(tx)
		ctx := context.Background()

		project := &domain.Project{WorkspaceID: workspace.ID, Name: "Launch"}
		if err := repo.Create(ctx, project); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if err := repo.Delete(ctx, project.ID); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		if _, err := repo.GetByID(ctx, project.ID); !errors.Is(err, domain.ErrProjectNotFound) {
			t.Errorf("GetByID() after Delete() error = %v, want %v", err, domain.ErrProjectNotFound)
		}

		got, err := repo.ListByWorkspaceID(ctx, workspace.ID)
		if err != nil {
			t.Fatalf("ListByWorkspaceID() error = %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected a deleted project to be excluded from ListByWorkspaceID, got %d", len(got))
		}
	})
}

func TestProjectRepository_Delete_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewProjectRepository(tx)

		err := repo.Delete(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrProjectNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, domain.ErrProjectNotFound)
		}
	})
}

// TestProjectRepository_Delete_CascadesToTasks exercises the documented
// cascade rule (docs/database.md) directly against the tasks table via raw
// SQL — there's no Task domain/repository yet (see 000003_projects.up.sql's
// comment), so this is the only way to verify Delete's cascade actually
// works until that vertical slice lands.
func TestProjectRepository_Delete_CascadesToTasks(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace := newTestProjectWorkspace(t, tx)
		repo := NewProjectRepository(tx)
		ctx := context.Background()

		project := &domain.Project{WorkspaceID: workspace.ID, Name: "Launch"}
		if err := repo.Create(ctx, project); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		taskID := uuid.New()
		if err := tx.Exec(
			`INSERT INTO tasks (id, workspace_id, project_id, title) VALUES (?, ?, ?, ?)`,
			taskID, workspace.ID, project.ID, "Ship it",
		).Error; err != nil {
			t.Fatalf("seeding task: %v", err)
		}

		if err := repo.Delete(ctx, project.ID); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		var deletedAt *string
		if err := tx.Raw(`SELECT deleted_at FROM tasks WHERE id = ?`, taskID).Scan(&deletedAt).Error; err != nil {
			t.Fatalf("querying task deleted_at: %v", err)
		}
		if deletedAt == nil {
			t.Error("expected the project's task to be soft-deleted (deleted_at set) after Delete()")
		}
	})
}
