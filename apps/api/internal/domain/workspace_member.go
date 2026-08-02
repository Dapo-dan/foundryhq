package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrWorkspaceMemberNotFound is returned by WorkspaceMemberRepository when no
// membership row matches the given lookup.
var ErrWorkspaceMemberNotFound = errors.New("workspace member not found")

// WorkspaceRole is one of the values allowed by the workspace_members.role
// CHECK constraint. Only Owner and Member are reachable from the API in v1 —
// Admin/Viewer are reserved in the schema but deferred to v1.1+ (docs/mvp.md).
type WorkspaceRole string

const (
	RoleOwner  WorkspaceRole = "owner"
	RoleAdmin  WorkspaceRole = "admin"
	RoleMember WorkspaceRole = "member"
	RoleViewer WorkspaceRole = "viewer"
)

// WorkspaceMember is a user's membership in a workspace, realizing the
// Users↔Workspaces many-to-many relationship plus a role.
type WorkspaceMember struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	UserID      uuid.UUID
	Role        WorkspaceRole
	InvitedAt   time.Time
	// JoinedAt is nil until the invite is accepted — set at creation time for
	// the owner membership WorkspaceUsecase.Create makes for the creator, and
	// by AuthUsecase.AcceptInvite for an invited member who activates their
	// placeholder account via the emailed invite token.
	JoinedAt *time.Time
}

// WorkspaceMemberWithUser augments WorkspaceMember with the member's email,
// for the members-list response (GET /workspaces/{id}/members) — the caller
// needs a human-readable identifier, not just a UserID.
type WorkspaceMemberWithUser struct {
	WorkspaceMember
	Email string
}

// WorkspaceMemberRepository persists and retrieves WorkspaceMember entities.
type WorkspaceMemberRepository interface {
	Create(ctx context.Context, member *WorkspaceMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*WorkspaceMember, error)
	GetByWorkspaceAndUser(ctx context.Context, workspaceID, userID uuid.UUID) (*WorkspaceMember, error)
	ListByWorkspaceIDWithUser(ctx context.Context, workspaceID uuid.UUID) ([]*WorkspaceMemberWithUser, error)
	UpdateRole(ctx context.Context, id uuid.UUID, role WorkspaceRole) error
	// MarkJoined sets JoinedAt to now for the membership row with the given
	// id — used by AuthUsecase.AcceptInvite once an invited user activates
	// their placeholder account.
	MarkJoined(ctx context.Context, id uuid.UUID) error
}
