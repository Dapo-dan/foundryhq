package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrSprintNotFound is returned by SprintRepository when no sprint matches
// the given lookup.
var ErrSprintNotFound = errors.New("sprint not found")

// Sprint is a time-boxed window a workspace plans tasks into. StartDate and
// EndDate are date-only (no time-of-day component) — see docs/database.md's
// Sprints section.
type Sprint struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SprintRepository persists and retrieves Sprint entities. There's no
// Update/Delete — docs/api.md doesn't call for either, and none of the v1
// UI needs them yet.
type SprintRepository interface {
	Create(ctx context.Context, sprint *Sprint) error
	GetByID(ctx context.Context, id uuid.UUID) (*Sprint, error)
	ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*Sprint, error)
}
