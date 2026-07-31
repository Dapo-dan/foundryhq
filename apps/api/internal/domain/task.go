package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrTaskNotFound is returned by TaskRepository when no task matches the
// given lookup.
var ErrTaskNotFound = errors.New("task not found")

// TaskStatus is one of the values allowed by the tasks.status CHECK
// constraint.
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

// TaskPriority is one of the values allowed by the tasks.priority CHECK
// constraint.
type TaskPriority string

const (
	PriorityUrgent TaskPriority = "urgent"
	PriorityHigh   TaskPriority = "high"
	PriorityMedium TaskPriority = "medium"
	PriorityLow    TaskPriority = "low"
)

// Task belongs to exactly one project and (denormalized) the workspace that
// project belongs to — see docs/database.md's Tasks section for why
// WorkspaceID is stored directly rather than derived through ProjectID.
type Task struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	ProjectID   uuid.UUID
	Title       string
	Status      TaskStatus
	// AssigneeID is nil for an unassigned task.
	AssigneeID *uuid.UUID
	// SprintID is nil for a backlog task (not in any sprint).
	SprintID    *uuid.UUID
	Priority    TaskPriority
	StoryPoints *int
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// DeletedAt is nil unless the task has been soft-deleted.
	DeletedAt *time.Time
}

// TaskFilter narrows ListByWorkspaceID's result. A nil field means
// "don't filter on this".
type TaskFilter struct {
	ProjectID  *uuid.UUID
	Status     *TaskStatus
	AssigneeID *uuid.UUID
	SprintID   *uuid.UUID
}

// TaskRepository persists and retrieves Task entities.
type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID, filter TaskFilter) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	// SumStoryPointsForSprint returns the sum of story_points for tasks in
	// sprintID with status = done and updated_at within [startDate, endDate]
	// (inclusive of the entire endDate calendar day) — the velocity
	// computation behind SprintUsecase.GetVelocity. Computed at query time,
	// never stored (docs/database.md's Sprints section).
	SumStoryPointsForSprint(ctx context.Context, sprintID uuid.UUID, startDate, endDate time.Time) (int, error)
}
