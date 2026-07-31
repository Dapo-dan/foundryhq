package postgres

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// These tests lock in migration 000006_fix_cascade_constraints — a
// hand-written DELETE via raw SQL (rather than through a repository method,
// since domain.WorkspaceRepository has no Delete and there's no v1 endpoint
// for it either) is the only way to actually trigger these FKs' ON DELETE
// behavior directly against Postgres.

func TestCascade_DeletingWorkspace_RemovesMembersAndTasks(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		ctx := context.Background()
		workspace, project := newTestTaskProject(t, tx)

		taskRepo := NewTaskRepository(tx)
		task := &domain.Task{WorkspaceID: workspace.ID, ProjectID: project.ID, Title: "Ship it"}
		if err := taskRepo.Create(ctx, task); err != nil {
			t.Fatalf("creating task: %v", err)
		}

		if err := tx.Exec(`DELETE FROM workspaces WHERE id = ?`, workspace.ID).Error; err != nil {
			t.Fatalf("deleting workspace: %v", err)
		}

		var memberCount, taskCount int64
		if err := tx.Raw(`SELECT count(*) FROM workspace_members WHERE workspace_id = ?`, workspace.ID).Scan(&memberCount).Error; err != nil {
			t.Fatalf("counting workspace_members: %v", err)
		}
		if memberCount != 0 {
			t.Errorf("workspace_members rows remaining = %d, want 0 — workspace_id should cascade", memberCount)
		}
		if err := tx.Raw(`SELECT count(*) FROM tasks WHERE workspace_id = ?`, workspace.ID).Scan(&taskCount).Error; err != nil {
			t.Fatalf("counting tasks: %v", err)
		}
		if taskCount != 0 {
			t.Errorf("tasks rows remaining = %d, want 0 — workspace_id should cascade", taskCount)
		}
	})
}

func TestCascade_DeletingUser_ClearsTaskAssignee(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		ctx := context.Background()
		workspace, project := newTestTaskProject(t, tx)

		userRepo := NewUserRepository(tx)
		assignee := newTestUser()
		if err := userRepo.Create(ctx, assignee); err != nil {
			t.Fatalf("creating assignee: %v", err)
		}

		taskRepo := NewTaskRepository(tx)
		task := &domain.Task{WorkspaceID: workspace.ID, ProjectID: project.ID, Title: "Ship it", AssigneeID: &assignee.ID}
		if err := taskRepo.Create(ctx, task); err != nil {
			t.Fatalf("creating task: %v", err)
		}

		if err := tx.Exec(`DELETE FROM users WHERE id = ?`, assignee.ID).Error; err != nil {
			t.Fatalf("deleting user: %v", err)
		}

		got, err := taskRepo.GetByID(ctx, task.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.AssigneeID != nil {
			t.Errorf("AssigneeID = %v, want nil after the assignee user was deleted", got.AssigneeID)
		}
	})
}
