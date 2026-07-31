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

// SprintUsecase implements sprint creation, lookup, and velocity
// computation, all scoped to a workspace. There's no Update/Delete —
// docs/api.md doesn't call for either.
type SprintUsecase struct {
	sprintRepo domain.SprintRepository
	taskRepo   domain.TaskRepository
	memberRepo domain.WorkspaceMemberRepository
}

// NewSprintUsecase constructs a SprintUsecase.
func NewSprintUsecase(sprintRepo domain.SprintRepository, taskRepo domain.TaskRepository, memberRepo domain.WorkspaceMemberRepository) *SprintUsecase {
	return &SprintUsecase{sprintRepo: sprintRepo, taskRepo: taskRepo, memberRepo: memberRepo}
}

// SprintWithTasks bundles a sprint with its tasks for GetByID's response.
// Grouping by status is a presentation concern handled by the caller (the
// handler's response shaping, or on web the same Kanban board component
// the Tasks page already uses) — not this usecase.
type SprintWithTasks struct {
	Sprint *domain.Sprint
	Tasks  []*domain.Task
}

// Create makes a new sprint in workspaceID, gated on callerID being a
// member of it. endDate must be on or after startDate — validated up front
// so a bad range fails with a clean field error rather than the database's
// raw CHECK-constraint violation.
func (u *SprintUsecase) Create(ctx context.Context, callerID, workspaceID uuid.UUID, name string, startDate, endDate time.Time) (*domain.Sprint, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperrors.Validation("name", "name is required")
	}
	if endDate.Before(startDate) {
		return nil, apperrors.Validation("endDate", "end date must be on or after the start date")
	}

	sprint := &domain.Sprint{WorkspaceID: workspaceID, Name: name, StartDate: startDate, EndDate: endDate}
	if err := u.sprintRepo.Create(ctx, sprint); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("creating sprint: %w", err))
	}
	return sprint, nil
}

// GetByID returns sprintID with its tasks, gated on callerID being a member
// of workspaceID and on sprintID actually belonging to workspaceID.
func (u *SprintUsecase) GetByID(ctx context.Context, callerID, workspaceID, sprintID uuid.UUID) (*SprintWithTasks, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	sprint, err := u.getScoped(ctx, workspaceID, sprintID)
	if err != nil {
		return nil, err
	}

	tasks, err := u.taskRepo.ListByWorkspaceID(ctx, workspaceID, domain.TaskFilter{SprintID: &sprintID})
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("listing sprint tasks: %w", err))
	}

	return &SprintWithTasks{Sprint: sprint, Tasks: tasks}, nil
}

// ListForWorkspace returns every sprint in workspaceID, gated on callerID
// being a member of it.
func (u *SprintUsecase) ListForWorkspace(ctx context.Context, callerID, workspaceID uuid.UUID) ([]*domain.Sprint, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	sprints, err := u.sprintRepo.ListByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("listing sprints: %w", err))
	}
	return sprints, nil
}

// GetVelocity returns the sum of story points completed within sprintID's
// date range (docs/database.md's Sprints section), computed live on every
// call — never cached or stored.
func (u *SprintUsecase) GetVelocity(ctx context.Context, callerID, workspaceID, sprintID uuid.UUID) (int, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return 0, err
	}

	sprint, err := u.getScoped(ctx, workspaceID, sprintID)
	if err != nil {
		return 0, err
	}

	velocity, err := u.taskRepo.SumStoryPointsForSprint(ctx, sprint.ID, sprint.StartDate, sprint.EndDate)
	if err != nil {
		return 0, apperrors.Internal(fmt.Errorf("computing velocity: %w", err))
	}
	return velocity, nil
}

func (u *SprintUsecase) requireMembership(ctx context.Context, workspaceID, userID uuid.UUID) error {
	_, err := u.memberRepo.GetByWorkspaceAndUser(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrWorkspaceMemberNotFound) {
			return apperrors.Forbidden("not a member of this workspace")
		}
		return apperrors.Internal(fmt.Errorf("checking workspace membership: %w", err))
	}
	return nil
}

// getScoped returns sprintID, but only if it actually belongs to
// workspaceID — otherwise it's treated as not found (same rationale as
// ProjectUsecase/TaskUsecase's getScoped).
func (u *SprintUsecase) getScoped(ctx context.Context, workspaceID, sprintID uuid.UUID) (*domain.Sprint, error) {
	sprint, err := u.sprintRepo.GetByID(ctx, sprintID)
	if err != nil {
		if errors.Is(err, domain.ErrSprintNotFound) {
			return nil, apperrors.NotFound("sprint not found")
		}
		return nil, apperrors.Internal(fmt.Errorf("getting sprint: %w", err))
	}
	if sprint.WorkspaceID != workspaceID {
		return nil, apperrors.NotFound("sprint not found")
	}
	return sprint, nil
}
