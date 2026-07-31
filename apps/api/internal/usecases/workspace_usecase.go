package usecases

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// maxSlugAttempts bounds the numeric-suffix retry generateUniqueSlug does
// when a workspace's default slug collides — a small, fixed cap rather than
// an unbounded loop against the database.
const maxSlugAttempts = 20

var slugInvalidChars = regexp.MustCompile(`[^a-z0-9]+`)

// WorkspaceUsecase implements workspace creation, lookup, updates, and team
// membership management.
type WorkspaceUsecase struct {
	workspaceRepo domain.WorkspaceRepository
	memberRepo    domain.WorkspaceMemberRepository
	userRepo      domain.UserRepository
}

// NewWorkspaceUsecase constructs a WorkspaceUsecase.
func NewWorkspaceUsecase(
	workspaceRepo domain.WorkspaceRepository,
	memberRepo domain.WorkspaceMemberRepository,
	userRepo domain.UserRepository,
) *WorkspaceUsecase {
	return &WorkspaceUsecase{workspaceRepo: workspaceRepo, memberRepo: memberRepo, userRepo: userRepo}
}

// UpdateWorkspaceInput carries the optional fields Update can change. A nil
// field is left unchanged — this is a partial-update (PATCH) input, not a
// full replacement.
type UpdateWorkspaceInput struct {
	Name    *string
	Slug    *string
	LogoURL *string
}

// Create makes a new workspace and, in the same transaction (see
// domain.WorkspaceRepository.Create), an owner membership for ownerID —
// every workspace always has exactly one owner (docs/database.md). The
// slug is generated server-side from name; callers never supply one.
func (u *WorkspaceUsecase) Create(ctx context.Context, ownerID uuid.UUID, name string) (*domain.Workspace, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperrors.Validation("name", "name is required")
	}

	slug, err := u.generateUniqueSlug(ctx, name)
	if err != nil {
		return nil, err
	}

	workspace := &domain.Workspace{Name: name, Slug: slug}
	now := time.Now()
	owner := &domain.WorkspaceMember{UserID: ownerID, Role: domain.RoleOwner, JoinedAt: &now}

	if err := u.workspaceRepo.Create(ctx, workspace, owner); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("creating workspace: %w", err))
	}
	return workspace, nil
}

// GetByID returns workspaceID, gated on callerID being a member of it.
func (u *WorkspaceUsecase) GetByID(ctx context.Context, callerID, workspaceID uuid.UUID) (*domain.Workspace, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}
	return u.getWorkspace(ctx, workspaceID)
}

// ListForUser returns every workspace userID belongs to. The web client
// treats an empty result as "hasn't finished onboarding yet".
func (u *WorkspaceUsecase) ListForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Workspace, error) {
	workspaces, err := u.workspaceRepo.ListForUser(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("listing workspaces: %w", err))
	}
	return workspaces, nil
}

// Update applies input's non-nil fields to workspaceID, gated on callerID
// being the workspace's owner — renaming/rebranding a shared workspace isn't
// a plain-member action. A Slug change is rejected with a conflict if it
// collides with another workspace's slug.
func (u *WorkspaceUsecase) Update(ctx context.Context, callerID, workspaceID uuid.UUID, input UpdateWorkspaceInput) (*domain.Workspace, error) {
	if err := u.requireOwner(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	workspace, err := u.getWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, apperrors.Validation("name", "name is required")
		}
		workspace.Name = name
	}
	if input.LogoURL != nil {
		workspace.LogoURL = strings.TrimSpace(*input.LogoURL)
	}
	if input.Slug != nil {
		slug := slugify(*input.Slug)
		if slug != workspace.Slug {
			exists, err := u.workspaceRepo.SlugExists(ctx, slug)
			if err != nil {
				return nil, apperrors.Internal(fmt.Errorf("checking slug uniqueness: %w", err))
			}
			if exists {
				return nil, apperrors.Conflict("slug already in use")
			}
			workspace.Slug = slug
		}
	}

	if err := u.workspaceRepo.Update(ctx, workspace); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("updating workspace: %w", err))
	}

	// Re-fetch rather than trust the in-memory copy: updated_at is set by the
	// database, not computed here.
	return u.getWorkspace(ctx, workspaceID)
}

// ListMembers returns every member of workspaceID (with email), gated on
// callerID being a member of it.
func (u *WorkspaceUsecase) ListMembers(ctx context.Context, callerID, workspaceID uuid.UUID) ([]*domain.WorkspaceMemberWithUser, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}
	members, err := u.memberRepo.ListByWorkspaceIDWithUser(ctx, workspaceID)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("listing workspace members: %w", err))
	}
	return members, nil
}

// Invite adds email to workspaceID as a member, gated on callerID being a
// member of it. If no user with that email exists yet, a placeholder User
// is created — email set, no password (domain.User.PasswordHash == "") —
// which AuthUsecase.Register later "claims" by setting a password the first
// time that address actually signs up, rather than rejecting it as a
// duplicate email.
func (u *WorkspaceUsecase) Invite(ctx context.Context, callerID, workspaceID uuid.UUID, email string) (*domain.WorkspaceMember, error) {
	if err := u.requireMembership(ctx, workspaceID, callerID); err != nil {
		return nil, err
	}

	email = normalizeEmail(email)
	if email == "" {
		return nil, apperrors.Validation("email", "email is required")
	}

	target, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return nil, apperrors.Internal(fmt.Errorf("looking up invitee: %w", err))
		}
		target = &domain.User{Email: email}
		if err := u.userRepo.Create(ctx, target); err != nil {
			return nil, apperrors.Internal(fmt.Errorf("creating placeholder user: %w", err))
		}
	} else {
		_, err := u.memberRepo.GetByWorkspaceAndUser(ctx, workspaceID, target.ID)
		if err == nil {
			return nil, apperrors.Conflict("this person is already a member")
		}
		if !errors.Is(err, domain.ErrWorkspaceMemberNotFound) {
			return nil, apperrors.Internal(fmt.Errorf("checking existing membership: %w", err))
		}
	}

	member := &domain.WorkspaceMember{WorkspaceID: workspaceID, UserID: target.ID, Role: domain.RoleMember}
	if err := u.memberRepo.Create(ctx, member); err != nil {
		return nil, apperrors.Internal(fmt.Errorf("creating workspace member: %w", err))
	}
	return member, nil
}

// UpdateMemberRole changes memberID's role within workspaceID, gated on
// callerID being the workspace's owner — a plain member must not be able to
// promote themselves (or anyone else) to owner, or demote the owner, which is
// why this checks ownership rather than plain membership. There's no
// transfer-ownership endpoint, so a change that would leave the workspace
// without an owner is refused with a conflict rather than performed.
func (u *WorkspaceUsecase) UpdateMemberRole(ctx context.Context, callerID, workspaceID, memberID uuid.UUID, role domain.WorkspaceRole) error {
	if err := u.requireOwner(ctx, workspaceID, callerID); err != nil {
		return err
	}
	if role != domain.RoleOwner && role != domain.RoleMember {
		return apperrors.Validation("role", "role must be one of: owner, member")
	}

	members, err := u.memberRepo.ListByWorkspaceIDWithUser(ctx, workspaceID)
	if err != nil {
		return apperrors.Internal(fmt.Errorf("listing workspace members: %w", err))
	}

	var target *domain.WorkspaceMemberWithUser
	ownerCount := 0
	for _, m := range members {
		if m.Role == domain.RoleOwner {
			ownerCount++
		}
		if m.ID == memberID {
			target = m
		}
	}
	if target == nil {
		return apperrors.NotFound("workspace member not found")
	}
	if target.Role == domain.RoleOwner && role != domain.RoleOwner && ownerCount <= 1 {
		return apperrors.Conflict("a workspace must always have an owner")
	}

	if err := u.memberRepo.UpdateRole(ctx, memberID, role); err != nil {
		return apperrors.Internal(fmt.Errorf("updating member role: %w", err))
	}
	return nil
}

// requireMembership returns apperrors.Forbidden unless userID is a member of
// workspaceID. Every WorkspaceUsecase method except Create starts with this
// guard — a workspace's data is only visible to its own members.
func (u *WorkspaceUsecase) requireMembership(ctx context.Context, workspaceID, userID uuid.UUID) error {
	_, err := u.memberRepo.GetByWorkspaceAndUser(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrWorkspaceMemberNotFound) {
			return apperrors.Forbidden("not a member of this workspace")
		}
		return apperrors.Internal(fmt.Errorf("checking workspace membership: %w", err))
	}
	return nil
}

// requireOwner returns apperrors.Forbidden unless userID is a member of
// workspaceID with the owner role. Used by the handful of actions
// (workspace settings, member role changes) that plain membership isn't
// enough to authorize — see docs/api.md's "forbidden" error, which is
// documented as "the workspace role doesn't permit this action", not just
// "not a member".
func (u *WorkspaceUsecase) requireOwner(ctx context.Context, workspaceID, userID uuid.UUID) error {
	member, err := u.memberRepo.GetByWorkspaceAndUser(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrWorkspaceMemberNotFound) {
			return apperrors.Forbidden("not a member of this workspace")
		}
		return apperrors.Internal(fmt.Errorf("checking workspace membership: %w", err))
	}
	if member.Role != domain.RoleOwner {
		return apperrors.Forbidden("only the workspace owner can perform this action")
	}
	return nil
}

func (u *WorkspaceUsecase) getWorkspace(ctx context.Context, workspaceID uuid.UUID) (*domain.Workspace, error) {
	workspace, err := u.workspaceRepo.GetByID(ctx, workspaceID)
	if err != nil {
		if errors.Is(err, domain.ErrWorkspaceNotFound) {
			return nil, apperrors.NotFound("workspace not found")
		}
		return nil, apperrors.Internal(fmt.Errorf("getting workspace: %w", err))
	}
	return workspace, nil
}

// generateUniqueSlug slugifies name and, on collision, retries with a
// numeric suffix up to maxSlugAttempts times.
func (u *WorkspaceUsecase) generateUniqueSlug(ctx context.Context, name string) (string, error) {
	base := slugify(name)
	slug := base
	for attempt := 0; attempt < maxSlugAttempts; attempt++ {
		if attempt > 0 {
			slug = fmt.Sprintf("%s-%d", base, attempt+1)
		}
		exists, err := u.workspaceRepo.SlugExists(ctx, slug)
		if err != nil {
			return "", apperrors.Internal(fmt.Errorf("checking slug uniqueness: %w", err))
		}
		if !exists {
			return slug, nil
		}
	}
	return "", apperrors.Internal(fmt.Errorf("could not generate a unique slug for %q after %d attempts", name, maxSlugAttempts))
}

// slugify lowercases name and collapses every run of non-alphanumeric
// characters into a single hyphen, trimming leading/trailing hyphens.
func slugify(name string) string {
	slug := slugInvalidChars.ReplaceAllString(strings.ToLower(name), "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "workspace"
	}
	return slug
}
