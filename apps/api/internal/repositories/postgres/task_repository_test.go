package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// newTestTaskProject creates a project (in a fresh workspace) that tests can
// hang tasks off of, returning both.
func newTestTaskProject(t *testing.T, tx *gorm.DB) (*domain.Workspace, *domain.Project) {
	t.Helper()
	workspace := newTestProjectWorkspace(t, tx)
	projectRepo := NewProjectRepository(tx)

	project := &domain.Project{WorkspaceID: workspace.ID, Name: "Launch"}
	if err := projectRepo.Create(context.Background(), project); err != nil {
		t.Fatalf("creating project: %v", err)
	}
	return workspace, project
}

func TestTaskRepository_CreateAndGetByID(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace, project := newTestTaskProject(t, tx)
		repo := NewTaskRepository(tx)
		ctx := context.Background()

		task := &domain.Task{WorkspaceID: workspace.ID, ProjectID: project.ID, Title: "Ship it"}
		if err := repo.Create(ctx, task); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if task.ID == uuid.Nil {
			t.Fatal("Create() did not populate ID")
		}
		if task.Status != domain.StatusTodo {
			t.Errorf("Status = %v, want default %v", task.Status, domain.StatusTodo)
		}

		got, err := repo.GetByID(ctx, task.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.Title != "Ship it" {
			t.Errorf("Title = %q, want %q", got.Title, "Ship it")
		}
	})
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewTaskRepository(tx)

		_, err := repo.GetByID(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("GetByID() error = %v, want %v", err, domain.ErrTaskNotFound)
		}
	})
}

func TestTaskRepository_ListByWorkspaceID_Filters(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace, projectA := newTestTaskProject(t, tx)
		projectRepo := NewProjectRepository(tx)
		userRepo := NewUserRepository(tx)
		repo := NewTaskRepository(tx)
		ctx := context.Background()

		projectB := &domain.Project{WorkspaceID: workspace.ID, Name: "Other"}
		if err := projectRepo.Create(ctx, projectB); err != nil {
			t.Fatalf("creating second project: %v", err)
		}

		assignee := newTestUser()
		if err := userRepo.Create(ctx, assignee); err != nil {
			t.Fatalf("creating assignee user: %v", err)
		}

		taskA := &domain.Task{WorkspaceID: workspace.ID, ProjectID: projectA.ID, Title: "In A", AssigneeID: &assignee.ID}
		if err := repo.Create(ctx, taskA); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		taskB := &domain.Task{WorkspaceID: workspace.ID, ProjectID: projectB.ID, Title: "In B"}
		if err := repo.Create(ctx, taskB); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		byProject, err := repo.ListByWorkspaceID(ctx, workspace.ID, domain.TaskFilter{ProjectID: &projectA.ID})
		if err != nil {
			t.Fatalf("ListByWorkspaceID() error = %v", err)
		}
		if len(byProject) != 1 || byProject[0].ID != taskA.ID {
			t.Errorf("filtering by projectId: got %d tasks, want just taskA", len(byProject))
		}

		byAssignee, err := repo.ListByWorkspaceID(ctx, workspace.ID, domain.TaskFilter{AssigneeID: &assignee.ID})
		if err != nil {
			t.Fatalf("ListByWorkspaceID() error = %v", err)
		}
		if len(byAssignee) != 1 || byAssignee[0].ID != taskA.ID {
			t.Errorf("filtering by assigneeId: got %d tasks, want just taskA", len(byAssignee))
		}

		todoStatus := domain.StatusTodo
		byStatus, err := repo.ListByWorkspaceID(ctx, workspace.ID, domain.TaskFilter{Status: &todoStatus})
		if err != nil {
			t.Fatalf("ListByWorkspaceID() error = %v", err)
		}
		if len(byStatus) != 2 {
			t.Errorf("filtering by status=todo: got %d tasks, want 2", len(byStatus))
		}
	})
}

func TestTaskRepository_Update(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace, project := newTestTaskProject(t, tx)
		repo := NewTaskRepository(tx)
		ctx := context.Background()

		task := &domain.Task{WorkspaceID: workspace.ID, ProjectID: project.ID, Title: "Ship it"}
		if err := repo.Create(ctx, task); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		task.Title = "Ship it faster"
		task.Status = domain.StatusInProgress
		if err := repo.Update(ctx, task); err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		got, err := repo.GetByID(ctx, task.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.Title != "Ship it faster" {
			t.Errorf("Title = %q, want %q", got.Title, "Ship it faster")
		}
		if got.Status != domain.StatusInProgress {
			t.Errorf("Status = %v, want %v", got.Status, domain.StatusInProgress)
		}
	})
}

func TestTaskRepository_Update_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewTaskRepository(tx)

		err := repo.Update(context.Background(), &domain.Task{ID: uuid.New(), Title: "x", Status: domain.StatusTodo})
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("Update() error = %v, want %v", err, domain.ErrTaskNotFound)
		}
	})
}

func TestTaskRepository_Delete_ExcludesFromFutureQueries(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace, project := newTestTaskProject(t, tx)
		repo := NewTaskRepository(tx)
		ctx := context.Background()

		task := &domain.Task{WorkspaceID: workspace.ID, ProjectID: project.ID, Title: "Ship it"}
		if err := repo.Create(ctx, task); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if err := repo.Delete(ctx, task.ID); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		if _, err := repo.GetByID(ctx, task.ID); !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("GetByID() after Delete() error = %v, want %v", err, domain.ErrTaskNotFound)
		}

		got, err := repo.ListByWorkspaceID(ctx, workspace.ID, domain.TaskFilter{})
		if err != nil {
			t.Fatalf("ListByWorkspaceID() error = %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected a deleted task to be excluded from ListByWorkspaceID, got %d", len(got))
		}
	})
}

func TestTaskRepository_Delete_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewTaskRepository(tx)

		err := repo.Delete(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, domain.ErrTaskNotFound)
		}
	})
}
