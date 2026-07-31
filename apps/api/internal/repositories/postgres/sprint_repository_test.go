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

func newTestSprint(workspaceID uuid.UUID) *domain.Sprint {
	return &domain.Sprint{
		WorkspaceID: workspaceID,
		Name:        "Sprint 1",
		StartDate:   time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC),
	}
}

func TestSprintRepository_CreateAndGetByID(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace := newTestProjectWorkspace(t, tx)
		repo := NewSprintRepository(tx)
		ctx := context.Background()

		sprint := newTestSprint(workspace.ID)
		if err := repo.Create(ctx, sprint); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if sprint.ID == uuid.Nil {
			t.Fatal("Create() did not populate ID")
		}

		got, err := repo.GetByID(ctx, sprint.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.Name != "Sprint 1" {
			t.Errorf("Name = %q, want %q", got.Name, "Sprint 1")
		}
		if !got.StartDate.Equal(sprint.StartDate) {
			t.Errorf("StartDate = %v, want %v", got.StartDate, sprint.StartDate)
		}
		if !got.EndDate.Equal(sprint.EndDate) {
			t.Errorf("EndDate = %v, want %v", got.EndDate, sprint.EndDate)
		}
	})
}

func TestSprintRepository_GetByID_NotFound(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		repo := NewSprintRepository(tx)

		_, err := repo.GetByID(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrSprintNotFound) {
			t.Errorf("GetByID() error = %v, want %v", err, domain.ErrSprintNotFound)
		}
	})
}

func TestSprintRepository_ListByWorkspaceID(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace := newTestProjectWorkspace(t, tx)
		otherWorkspace := newTestProjectWorkspace(t, tx)
		repo := NewSprintRepository(tx)
		ctx := context.Background()

		if err := repo.Create(ctx, newTestSprint(workspace.ID)); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if err := repo.Create(ctx, newTestSprint(otherWorkspace.ID)); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.ListByWorkspaceID(ctx, workspace.ID)
		if err != nil {
			t.Fatalf("ListByWorkspaceID() error = %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 sprint, got %d", len(got))
		}
		if got[0].WorkspaceID != workspace.ID {
			t.Errorf("WorkspaceID = %v, want %v", got[0].WorkspaceID, workspace.ID)
		}
	})
}

// TestSprintRepository_EndDateBeforeStartDateViolatesCheck confirms the
// DB-level CHECK constraint is actually in place as a backstop — the
// usecase validates this up front, but the constraint should still reject
// a bad range if something ever bypasses the usecase.
func TestSprintRepository_EndDateBeforeStartDateViolatesCheck(t *testing.T) {
	withTestTx(t, func(tx *gorm.DB) {
		workspace := newTestProjectWorkspace(t, tx)
		repo := NewSprintRepository(tx)

		sprint := &domain.Sprint{
			WorkspaceID: workspace.ID,
			Name:        "Backwards sprint",
			StartDate:   time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		}
		if err := repo.Create(context.Background(), sprint); err == nil {
			t.Error("Create() with end_date before start_date should fail the DB CHECK constraint")
		}
	})
}
