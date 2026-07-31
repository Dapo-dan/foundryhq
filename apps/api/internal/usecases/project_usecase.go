package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// ProjectUsecase implements project creation, lookup, updates, and deletion,
// all scoped to a workspace.
type ProjectUsecase struct {
	projectRepo domain.ProjectRepository
	memberRepo  domain.WorkspaceMemberRepository
}

// NewProjectUsecase constructs a ProjectUsecase.
func NewProjectUsecase(projectRepo domain.ProjectRepository, memberRepo domain.WorkspaceMemberRepository) *ProjectUsecase {
	return &ProjectUsecase{projectRepo: projectRepo, memberRepo: memberRepo}
}

// UpdateProjectInput carries the optional fields Update can change. A nil
// field is left unchanged — this is a partial-update (PATCH) input.
type UpdateProjectInput struct {
	Name        *string
	Description *string
}

// Create makes a new project in workspaceID, gated on callerID being a
// member of it.
func (u *ProjectUsecase) Create(ctx context.Context, callerID, workspaceID uuid.UUID, name string, description *string) (*domain.Project, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperrors.Validation("name", "name is required")
	}

	project := &domain.Project{WorkspaceID: workspaceID, Name: name, Description: trimmedOrNil(description)}
	if err := u.projectRepo.Create(ctx, project); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("creating project: %w", err))
	}
	return project, nil
}

// GetByID returns projectID, gated on callerID being a member of
// workspaceID and on projectID actually belonging to workspaceID.
func (u *ProjectUsecase) GetByID(ctx context.Context, callerID, workspaceID, projectID uuid.UUID) (*domain.Project, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}
	return u.getScoped(ctx, workspaceID, projectID)
}

// ListByWorkspaceID returns every project in workspaceID, gated on callerID
// being a member of it.
func (u *ProjectUsecase) ListByWorkspaceID(ctx context.Context, callerID, workspaceID uuid.UUID) ([]*domain.Project, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	projects, err := u.projectRepo.ListByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("listing projects: %w", err))
	}
	return projects, nil
}

// Update applies input's non-nil fields to projectID, gated the same way as
// GetByID.
func (u *ProjectUsecase) Update(ctx context.Context, callerID, workspaceID, projectID uuid.UUID, input UpdateProjectInput) (*domain.Project, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	project, err := u.getScoped(ctx, workspaceID, projectID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, apperrors.Validation("name", "name is required")
		}
		project.Name = name
	}
	if input.Description != nil {
		project.Description = trimmedOrNil(input.Description)
	}

	if err := u.projectRepo.Update(ctx, project); err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			return nil, apperrors.NotFound("project not found")
		}
		return nil, apperrors.Internal(fmt.Errorf("updating project: %w", err))
	}

	// Re-fetch rather than trust the in-memory copy: updated_at is set by the
	// database, not computed here.
	return u.getScoped(ctx, workspaceID, projectID)
}

// Delete soft-deletes projectID and its tasks, gated the same way as
// GetByID.
func (u *ProjectUsecase) Delete(ctx context.Context, callerID, workspaceID, projectID uuid.UUID) error {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return err
	}
	if _, err := u.getScoped(ctx, workspaceID, projectID); err != nil {
		return err
	}

	if err := u.projectRepo.Delete(ctx, projectID); err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			return apperrors.NotFound("project not found")
		}
		return apperrors.Internal(fmt.Errorf("deleting project: %w", err))
	}
	return nil
}

// requireMembership returns apperrors.Forbidden unless userID is a member of
// workspaceID.
func (u *ProjectUsecase) requireMembership(ctx context.Context, workspaceID, userID uuid.UUID) error {
	_, err := u.memberRepo.GetByWorkspaceAndUser(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrWorkspaceMemberNotFound) {
			return apperrors.Forbidden("not a member of this workspace")
		}
		return apperrors.Internal(fmt.Errorf("checking workspace membership: %w", err))
	}
	return nil
}

// getScoped returns projectID, but only if it actually belongs to
// workspaceID — otherwise it's treated as not found, exactly like a project
// that doesn't exist at all. requireMembership alone only confirms callerID
// belongs to workspaceID; it says nothing about which workspace projectID is
// actually in, so without this check a caller could read/mutate another
// workspace's project just by guessing its UUID (see docs/api.md's
// not_found row: "exists in a different workspace").
func (u *ProjectUsecase) getScoped(ctx context.Context, workspaceID, projectID uuid.UUID) (*domain.Project, error) {
	project, err := u.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			return nil, apperrors.NotFound("project not found")
		}
		return nil, apperrors.Internal(fmt.Errorf("getting project: %w", err))
	}
	if project.WorkspaceID != workspaceID {
		return nil, apperrors.NotFound("project not found")
	}
	return project, nil
}

// trimmedOrNil trims s and returns nil if either s or the trimmed result is
// empty — used for Description, which is optional (nil), not required-but-
// blankable.
func trimmedOrNil(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
