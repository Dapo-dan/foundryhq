package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// TaskUsecase implements task creation, lookup, updates, and deletion, all
// scoped to a workspace.
type TaskUsecase struct {
	taskRepo    domain.TaskRepository
	projectRepo domain.ProjectRepository
	sprintRepo  domain.SprintRepository
	memberRepo  domain.WorkspaceMemberRepository
}

// NewTaskUsecase constructs a TaskUsecase.
func NewTaskUsecase(taskRepo domain.TaskRepository, projectRepo domain.ProjectRepository, sprintRepo domain.SprintRepository, memberRepo domain.WorkspaceMemberRepository) *TaskUsecase {
	return &TaskUsecase{taskRepo: taskRepo, projectRepo: projectRepo, sprintRepo: sprintRepo, memberRepo: memberRepo}
}

// CreateTaskInput carries Create's required and optional fields.
type CreateTaskInput struct {
	ProjectID   uuid.UUID
	Title       string
	AssigneeID  *uuid.UUID
	SprintID    *uuid.UUID
	Priority    *domain.TaskPriority
	StoryPoints *int
	DueDate     *time.Time
}

// UpdateTaskInput carries the optional fields Update can change. A nil
// field is left unchanged — this is a partial-update (PATCH) input.
type UpdateTaskInput struct {
	ProjectID  *uuid.UUID
	Title      *string
	Status     *domain.TaskStatus
	AssigneeID *uuid.UUID
	// ClearAssignee unassigns the task regardless of AssigneeID — needed
	// because AssigneeID's own nil already means "leave unchanged", so a
	// plain nil can't also mean "unassign".
	ClearAssignee bool
	SprintID *uuid.UUID
	// ClearSprint moves the task back to the backlog regardless of
	// SprintID — same rationale as ClearAssignee (SprintID's own nil already
	// means "leave unchanged", so clearing needs its own explicit signal).
	ClearSprint bool
	Priority    *domain.TaskPriority
	StoryPoints *int
	// ClearStoryPoints clears story points regardless of StoryPoints — same
	// rationale as ClearAssignee (a plain nil already means "leave
	// unchanged", so clearing needs its own explicit signal).
	ClearStoryPoints bool
	DueDate          *time.Time
	// ClearDueDate clears the due date regardless of DueDate — same
	// rationale as ClearAssignee.
	ClearDueDate bool
}

// Create makes a new task in workspaceID, gated on callerID being a member
// of it. ProjectID must belong to workspaceID, and AssigneeID (if set) must
// be a member of workspaceID — both enforced here, not at the database
// level (docs/database.md's Tasks section: "enforced in the usecase layer
// when a task is created or its project_id changes").
func (u *TaskUsecase) Create(ctx context.Context, callerID, workspaceID uuid.UUID, input CreateTaskInput) (*domain.Task, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, apperrors.Validation("title", "title is required")
	}

	if err := u.requireProjectInWorkspace(ctx, workspaceID, input.ProjectID); err != nil {
		return nil, err
	}
	if input.AssigneeID != nil {
		if err := u.validateAssignee(ctx, workspaceID, *input.AssigneeID); err != nil {
			return nil, err
		}
	}
	if input.SprintID != nil {
		if err := u.requireSprintInWorkspace(ctx, workspaceID, *input.SprintID); err != nil {
			return nil, err
		}
	}
	priority := domain.PriorityMedium
	if input.Priority != nil {
		if !isValidTaskPriority(*input.Priority) {
			return nil, apperrors.Validation("priority", "priority must be one of: urgent, high, medium, low")
		}
		priority = *input.Priority
	}
	if input.StoryPoints != nil && *input.StoryPoints < 0 {
		return nil, apperrors.Validation("storyPoints", "story points must be zero or greater")
	}

	task := &domain.Task{
		WorkspaceID: workspaceID,
		ProjectID:   input.ProjectID,
		Title:       title,
		Status:      domain.StatusTodo,
		AssigneeID:  input.AssigneeID,
		SprintID:    input.SprintID,
		Priority:    priority,
		StoryPoints: input.StoryPoints,
		DueDate:     input.DueDate,
	}
	if err := u.taskRepo.Create(ctx, task); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("creating task: %w", err))
	}
	return task, nil
}

// GetByID returns taskID, gated on callerID being a member of workspaceID
// and on taskID actually belonging to workspaceID.
func (u *TaskUsecase) GetByID(ctx context.Context, callerID, workspaceID, taskID uuid.UUID) (*domain.Task, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}
	return u.getScoped(ctx, workspaceID, taskID)
}

// ListByWorkspaceID returns every task in workspaceID matching filter,
// gated on callerID being a member of it.
func (u *TaskUsecase) ListByWorkspaceID(ctx context.Context, callerID, workspaceID uuid.UUID, filter domain.TaskFilter) ([]*domain.Task, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	tasks, err := u.taskRepo.ListByWorkspaceID(ctx, workspaceID, filter)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("listing tasks: %w", err))
	}
	return tasks, nil
}

// Update applies input's non-nil fields to taskID, gated the same way as
// GetByID. A ProjectID change is validated exactly like Create's — the new
// project must belong to workspaceID too.
func (u *TaskUsecase) Update(ctx context.Context, callerID, workspaceID, taskID uuid.UUID, input UpdateTaskInput) (*domain.Task, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	task, err := u.getScoped(ctx, workspaceID, taskID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, apperrors.Validation("title", "title is required")
		}
		task.Title = title
	}
	if input.ProjectID != nil && *input.ProjectID != task.ProjectID {
		if err := u.requireProjectInWorkspace(ctx, workspaceID, *input.ProjectID); err != nil {
			return nil, err
		}
		task.ProjectID = *input.ProjectID
	}
	if input.Status != nil {
		if !isValidTaskStatus(*input.Status) {
			return nil, apperrors.Validation("status", "status must be one of: todo, in_progress, done")
		}
		task.Status = *input.Status
	}
	if input.ClearAssignee {
		task.AssigneeID = nil
	} else if input.AssigneeID != nil {
		if err := u.validateAssignee(ctx, workspaceID, *input.AssigneeID); err != nil {
			return nil, err
		}
		task.AssigneeID = input.AssigneeID
	}
	if input.ClearSprint {
		task.SprintID = nil
	} else if input.SprintID != nil && (task.SprintID == nil || *input.SprintID != *task.SprintID) {
		if err := u.requireSprintInWorkspace(ctx, workspaceID, *input.SprintID); err != nil {
			return nil, err
		}
		task.SprintID = input.SprintID
	}
	if input.Priority != nil {
		if !isValidTaskPriority(*input.Priority) {
			return nil, apperrors.Validation("priority", "priority must be one of: urgent, high, medium, low")
		}
		task.Priority = *input.Priority
	}
	if input.ClearStoryPoints {
		task.StoryPoints = nil
	} else if input.StoryPoints != nil {
		if *input.StoryPoints < 0 {
			return nil, apperrors.Validation("storyPoints", "story points must be zero or greater")
		}
		task.StoryPoints = input.StoryPoints
	}
	if input.ClearDueDate {
		task.DueDate = nil
	} else if input.DueDate != nil {
		task.DueDate = input.DueDate
	}

	if err := u.taskRepo.Update(ctx, task); err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, apperrors.NotFound("task not found")
		}
		return nil, apperrors.Internal(fmt.Errorf("updating task: %w", err))
	}

	// Re-fetch rather than trust the in-memory copy: updated_at is set by the
	// database, not computed here.
	return u.getScoped(ctx, workspaceID, taskID)
}

// Delete soft-deletes taskID, gated the same way as GetByID.
func (u *TaskUsecase) Delete(ctx context.Context, callerID, workspaceID, taskID uuid.UUID) error {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return err
	}
	if _, err := u.getScoped(ctx, workspaceID, taskID); err != nil {
		return err
	}

	if err := u.taskRepo.Delete(ctx, taskID); err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return apperrors.NotFound("task not found")
		}
		return apperrors.Internal(fmt.Errorf("deleting task: %w", err))
	}
	return nil
}

// isMember reports whether userID belongs to workspaceID, distinguishing
// "not a member" (false, nil) from an actual lookup failure (false, err).
func (u *TaskUsecase) isMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error) {
	_, err := u.memberRepo.GetByWorkspaceAndUser(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrWorkspaceMemberNotFound) {
			return false, nil
		}
		return false, apperrors.Internal(fmt.Errorf("checking workspace membership: %w", err))
	}
	return true, nil
}

// requireMembership returns apperrors.Forbidden unless userID is a member of
// workspaceID.
func (u *TaskUsecase) requireMembership(ctx context.Context, workspaceID, userID uuid.UUID) error {
	ok, err := u.isMember(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return apperrors.Forbidden("not a member of this workspace")
	}
	return nil
}

// validateAssignee returns apperrors.Validation unless assigneeID is a
// member of workspaceID — a task can't be assigned to someone outside the
// workspace it belongs to.
func (u *TaskUsecase) validateAssignee(ctx context.Context, workspaceID, assigneeID uuid.UUID) error {
	ok, err := u.isMember(ctx, workspaceID, assigneeID)
	if err != nil {
		return err
	}
	if !ok {
		return apperrors.Validation("assigneeId", "assignee is not a member of this workspace")
	}
	return nil
}

// requireProjectInWorkspace returns apperrors.Validation unless projectID
// exists and belongs to workspaceID — see docs/database.md's
// tasks.workspace_id == projects.workspace_id invariant.
func (u *TaskUsecase) requireProjectInWorkspace(ctx context.Context, workspaceID, projectID uuid.UUID) error {
	project, err := u.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			return apperrors.Validation("projectId", "project not found")
		}
		return apperrors.Internal(fmt.Errorf("getting project: %w", err))
	}
	if project.WorkspaceID != workspaceID {
		return apperrors.Validation("projectId", "project does not belong to this workspace")
	}
	return nil
}

// requireSprintInWorkspace returns apperrors.Validation unless sprintID
// exists and belongs to workspaceID — same invariant as
// requireProjectInWorkspace, applied to sprint assignment.
func (u *TaskUsecase) requireSprintInWorkspace(ctx context.Context, workspaceID, sprintID uuid.UUID) error {
	sprint, err := u.sprintRepo.GetByID(ctx, sprintID)
	if err != nil {
		if errors.Is(err, domain.ErrSprintNotFound) {
			return apperrors.Validation("sprintId", "sprint not found")
		}
		return apperrors.Internal(fmt.Errorf("getting sprint: %w", err))
	}
	if sprint.WorkspaceID != workspaceID {
		return apperrors.Validation("sprintId", "sprint does not belong to this workspace")
	}
	return nil
}

// getScoped returns taskID, but only if it actually belongs to
// workspaceID — otherwise it's treated as not found, exactly like a task
// that doesn't exist at all (same rationale as ProjectUsecase.getScoped).
func (u *TaskUsecase) getScoped(ctx context.Context, workspaceID, taskID uuid.UUID) (*domain.Task, error) {
	task, err := u.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, apperrors.NotFound("task not found")
		}
		return nil, apperrors.Internal(fmt.Errorf("getting task: %w", err))
	}
	if task.WorkspaceID != workspaceID {
		return nil, apperrors.NotFound("task not found")
	}
	return task, nil
}

func isValidTaskStatus(s domain.TaskStatus) bool {
	switch s {
	case domain.StatusTodo, domain.StatusInProgress, domain.StatusDone:
		return true
	default:
		return false
	}
}

func isValidTaskPriority(p domain.TaskPriority) bool {
	switch p {
	case domain.PriorityUrgent, domain.PriorityHigh, domain.PriorityMedium, domain.PriorityLow:
		return true
	default:
		return false
	}
}
