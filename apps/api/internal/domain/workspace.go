package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrWorkspaceNotFound is returned by WorkspaceRepository when no workspace
// matches the given lookup.
var ErrWorkspaceNotFound = errors.New("workspace not found")

// Workspace is a tenant boundary — every task, project, and team membership
// belongs to exactly one (see docs/database.md's Workspaces section).
type Workspace struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	LogoURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WorkspaceRepository persists and retrieves Workspace entities.
type WorkspaceRepository interface {
	// Create inserts workspace and owner's membership row together in one
	// transaction — a workspace must never exist without exactly one owner
	// (docs/database.md), so the two inserts succeed or fail as a unit.
	Create(ctx context.Context, workspace *Workspace, owner *WorkspaceMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*Workspace, error)
	Update(ctx context.Context, workspace *Workspace) error
	SlugExists(ctx context.Context, slug string) (bool, error)
	// ListForUser returns every workspace userID is a member of.
	ListForUser(ctx context.Context, userID uuid.UUID) ([]*Workspace, error)
}
