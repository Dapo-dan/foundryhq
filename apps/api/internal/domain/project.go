package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrProjectNotFound is returned by ProjectRepository when no project
// matches the given lookup.
var ErrProjectNotFound = errors.New("project not found")

// Project is the organizational unit tasks live under (docs/database.md) —
// a workspace can have many; every task belongs to exactly one.
type Project struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// DeletedAt is nil unless the project has been soft-deleted.
	DeletedAt *time.Time
}

// ProjectRepository persists and retrieves Project entities.
type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*Project, error)
	Update(ctx context.Context, project *Project) error
	// Delete soft-deletes project and, in the same operation, every task
	// belonging to it — Postgres has no way to cascade a soft-delete
	// (deleted_at is a column, not a row removal), so implementations must
	// enforce this explicitly (docs/database.md's Projects cascade rule).
	Delete(ctx context.Context, id uuid.UUID) error
}
