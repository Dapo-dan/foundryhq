package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// mockWorkspaceMemberRepo is a hand-written in-memory
// domain.WorkspaceMemberRepository, per CONTRIBUTING.md's "mock
// repositories via interfaces" convention.
type mockWorkspaceMemberRepo struct {
	byID map[uuid.UUID]*domain.WorkspaceMember
	// emails backs the "with user" half of ListByWorkspaceIDWithUser — the
	// real repository joins to the users table for this, so tests populate
	// it directly instead of wiring up a mockUserRepo dependency.
	emails map[uuid.UUID]string
}

func newMockWorkspaceMemberRepo() *mockWorkspaceMemberRepo {
	return &mockWorkspaceMemberRepo{
		byID:   map[uuid.UUID]*domain.WorkspaceMember{},
		emails: map[uuid.UUID]string{},
	}
}

func (m *mockWorkspaceMemberRepo) Create(_ context.Context, member *domain.WorkspaceMember) error {
	member.ID = uuid.New()
	member.InvitedAt = time.Now()
	cp := *member
	m.byID[member.ID] = &cp
	return nil
}

func (m *mockWorkspaceMemberRepo) GetByWorkspaceAndUser(_ context.Context, workspaceID, userID uuid.UUID) (*domain.WorkspaceMember, error) {
	for _, member := range m.byID {
		if member.WorkspaceID == workspaceID && member.UserID == userID {
			cp := *member
			return &cp, nil
		}
	}
	return nil, domain.ErrWorkspaceMemberNotFound
}

func (m *mockWorkspaceMemberRepo) ListByWorkspaceIDWithUser(_ context.Context, workspaceID uuid.UUID) ([]*domain.WorkspaceMemberWithUser, error) {
	var result []*domain.WorkspaceMemberWithUser
	for _, member := range m.byID {
		if member.WorkspaceID == workspaceID {
			result = append(result, &domain.WorkspaceMemberWithUser{
				WorkspaceMember: *member,
				Email:           m.emails[member.UserID],
			})
		}
	}
	return result, nil
}

func (m *mockWorkspaceMemberRepo) UpdateRole(_ context.Context, id uuid.UUID, role domain.WorkspaceRole) error {
	member, ok := m.byID[id]
	if !ok {
		return domain.ErrWorkspaceMemberNotFound
	}
	member.Role = role
	return nil
}

// mockWorkspaceRepo is a hand-written in-memory domain.WorkspaceRepository.
// Its Create writes the owner membership through the shared members store
// (rather than into its own private state), mirroring the real Postgres
// repository's single-transaction insert of both rows.
type mockWorkspaceRepo struct {
	byID    map[uuid.UUID]*domain.Workspace
	bySlug  map[string]bool
	members *mockWorkspaceMemberRepo
}

func newMockWorkspaceRepo(members *mockWorkspaceMemberRepo) *mockWorkspaceRepo {
	return &mockWorkspaceRepo{
		byID:    map[uuid.UUID]*domain.Workspace{},
		bySlug:  map[string]bool{},
		members: members,
	}
}

func (m *mockWorkspaceRepo) Create(ctx context.Context, workspace *domain.Workspace, owner *domain.WorkspaceMember) error {
	workspace.ID = uuid.New()
	workspace.CreatedAt = time.Now()
	workspace.UpdatedAt = workspace.CreatedAt
	cp := *workspace
	m.byID[workspace.ID] = &cp
	m.bySlug[workspace.Slug] = true

	owner.WorkspaceID = workspace.ID
	return m.members.Create(ctx, owner)
}

func (m *mockWorkspaceRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Workspace, error) {
	w, ok := m.byID[id]
	if !ok {
		return nil, domain.ErrWorkspaceNotFound
	}
	cp := *w
	return &cp, nil
}

func (m *mockWorkspaceRepo) Update(_ context.Context, workspace *domain.Workspace) error {
	if _, ok := m.byID[workspace.ID]; !ok {
		return domain.ErrWorkspaceNotFound
	}
	cp := *workspace
	cp.UpdatedAt = time.Now()
	m.byID[workspace.ID] = &cp
	m.bySlug[cp.Slug] = true
	return nil
}

func (m *mockWorkspaceRepo) SlugExists(_ context.Context, slug string) (bool, error) {
	return m.bySlug[slug], nil
}

func (m *mockWorkspaceRepo) ListForUser(_ context.Context, userID uuid.UUID) ([]*domain.Workspace, error) {
	var result []*domain.Workspace
	for _, member := range m.members.byID {
		if member.UserID == userID {
			if w, ok := m.byID[member.WorkspaceID]; ok {
				cp := *w
				result = append(result, &cp)
			}
		}
	}
	return result, nil
}

func newTestWorkspaceUsecase() (*WorkspaceUsecase, *mockWorkspaceRepo, *mockWorkspaceMemberRepo, *mockUserRepo) {
	users := newMockUserRepo()
	members := newMockWorkspaceMemberRepo()
	workspaces := newMockWorkspaceRepo(members)
	u := NewWorkspaceUsecase(workspaces, members, users)
	return u, workspaces, members, users
}

func TestWorkspaceCreate_Success(t *testing.T) {
	u, _, members, _ := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if workspace.Slug != "acme-inc" {
		t.Errorf("Slug = %q, want %q", workspace.Slug, "acme-inc")
	}

	member, err := members.GetByWorkspaceAndUser(context.Background(), workspace.ID, ownerID)
	if err != nil {
		t.Fatalf("expected an owner membership to be created, got error = %v", err)
	}
	if member.Role != domain.RoleOwner {
		t.Errorf("Role = %v, want %v", member.Role, domain.RoleOwner)
	}
	if member.JoinedAt == nil {
		t.Error("owner's JoinedAt should be set immediately, not nil")
	}
}

func TestWorkspaceCreate_EmptyName(t *testing.T) {
	u, _, _, _ := newTestWorkspaceUsecase()

	_, err := u.Create(context.Background(), uuid.New(), "   ")
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
}

func TestWorkspaceCreate_SlugCollisionAppendsSuffix(t *testing.T) {
	u, workspaces, _, _ := newTestWorkspaceUsecase()
	workspaces.bySlug["acme-inc"] = true

	workspace, err := u.Create(context.Background(), uuid.New(), "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if workspace.Slug == "acme-inc" {
		t.Error("Create() should have retried with a numeric suffix on slug collision")
	}
}

func TestGetByID_ForbiddenForNonMember(t *testing.T) {
	u, _, _, _ := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = u.GetByID(context.Background(), uuid.New(), workspace.ID)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeForbidden {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeForbidden)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	u, _, members, _ := newTestWorkspaceUsecase()
	callerID := uuid.New()
	workspaceID := uuid.New()

	// Membership exists but the workspace row itself doesn't — exercises the
	// GetByID-after-the-membership-guard NotFound path specifically.
	if err := members.Create(context.Background(), &domain.WorkspaceMember{WorkspaceID: workspaceID, UserID: callerID, Role: domain.RoleOwner}); err != nil {
		t.Fatalf("seeding membership: %v", err)
	}

	_, err := u.GetByID(context.Background(), callerID, workspaceID)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeNotFound)
	}
}

func TestUpdate_Success(t *testing.T) {
	u, _, _, _ := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	newName := "Acme Corp"
	newLogo := "https://example.com/logo.png"
	updated, err := u.Update(context.Background(), ownerID, workspace.ID, UpdateWorkspaceInput{
		Name:    &newName,
		LogoURL: &newLogo,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Name != newName {
		t.Errorf("Name = %q, want %q", updated.Name, newName)
	}
	if updated.LogoURL != newLogo {
		t.Errorf("LogoURL = %q, want %q", updated.LogoURL, newLogo)
	}
	if updated.Slug != workspace.Slug {
		t.Errorf("Slug = %q, want unchanged %q", updated.Slug, workspace.Slug)
	}
}

func TestUpdate_SlugConflict(t *testing.T) {
	u, _, _, _ := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	if _, err := u.Create(context.Background(), uuid.New(), "Other Workspace"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	takenSlug := "other-workspace"
	_, err = u.Update(context.Background(), ownerID, workspace.ID, UpdateWorkspaceInput{Slug: &takenSlug})
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeConflict {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeConflict)
	}
}

func TestListForUser(t *testing.T) {
	u, _, _, _ := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	if _, err := u.Create(context.Background(), ownerID, "Acme Inc."); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := u.Create(context.Background(), uuid.New(), "Someone Else's Workspace"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	workspaces, err := u.ListForUser(context.Background(), ownerID)
	if err != nil {
		t.Fatalf("ListForUser() error = %v", err)
	}
	if len(workspaces) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(workspaces))
	}
	if workspaces[0].Name != "Acme Inc." {
		t.Errorf("Name = %q, want %q", workspaces[0].Name, "Acme Inc.")
	}
}

func TestListMembers_ForbiddenForNonMember(t *testing.T) {
	u, _, _, _ := newTestWorkspaceUsecase()
	workspace, err := u.Create(context.Background(), uuid.New(), "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = u.ListMembers(context.Background(), uuid.New(), workspace.ID)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeForbidden {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeForbidden)
	}
}

func TestInvite_NewEmailCreatesPlaceholderUser(t *testing.T) {
	u, _, members, users := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	member, err := u.Invite(context.Background(), ownerID, workspace.ID, "New@Example.com")
	if err != nil {
		t.Fatalf("Invite() error = %v", err)
	}
	if member.Role != domain.RoleMember {
		t.Errorf("Role = %v, want %v", member.Role, domain.RoleMember)
	}
	if member.JoinedAt != nil {
		t.Error("an invited-but-not-yet-accepted member should have a nil JoinedAt")
	}

	placeholder, err := users.GetByEmail(context.Background(), "new@example.com")
	if err != nil {
		t.Fatalf("expected a placeholder user to be created, got error = %v", err)
	}
	if placeholder.PasswordHash != "" {
		t.Error("placeholder user should have no password set")
	}

	if _, err := members.GetByWorkspaceAndUser(context.Background(), workspace.ID, placeholder.ID); err != nil {
		t.Errorf("expected a membership row for the placeholder user, got error = %v", err)
	}
}

func TestInvite_ExistingUserBecomesMember(t *testing.T) {
	u, _, _, users := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	existing := &domain.User{Email: "existing@example.com", PasswordHash: "hashed"}
	if err := users.Create(context.Background(), existing); err != nil {
		t.Fatalf("seeding existing user: %v", err)
	}

	member, err := u.Invite(context.Background(), ownerID, workspace.ID, "existing@example.com")
	if err != nil {
		t.Fatalf("Invite() error = %v", err)
	}
	if member.UserID != existing.ID {
		t.Errorf("UserID = %v, want %v", member.UserID, existing.ID)
	}
}

func TestInvite_AlreadyMemberIsConflict(t *testing.T) {
	u, _, _, users := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	existing := &domain.User{Email: "existing@example.com", PasswordHash: "hashed"}
	if err := users.Create(context.Background(), existing); err != nil {
		t.Fatalf("seeding existing user: %v", err)
	}
	if _, err := u.Invite(context.Background(), ownerID, workspace.ID, "existing@example.com"); err != nil {
		t.Fatalf("first Invite() error = %v", err)
	}

	_, err = u.Invite(context.Background(), ownerID, workspace.ID, "existing@example.com")
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeConflict {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeConflict)
	}
}

func TestUpdateMemberRole_Success(t *testing.T) {
	u, _, members, users := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	newOwner := &domain.User{Email: "member@example.com", PasswordHash: "hashed"}
	if err := users.Create(context.Background(), newOwner); err != nil {
		t.Fatalf("seeding user: %v", err)
	}
	invited, err := u.Invite(context.Background(), ownerID, workspace.ID, "member@example.com")
	if err != nil {
		t.Fatalf("Invite() error = %v", err)
	}

	// Promote the invited member to owner — the workspace now has two
	// owners, so demoting the original owner afterward must succeed.
	if err := u.UpdateMemberRole(context.Background(), ownerID, workspace.ID, invited.ID, domain.RoleOwner); err != nil {
		t.Fatalf("UpdateMemberRole() promote error = %v", err)
	}

	originalOwner, err := members.GetByWorkspaceAndUser(context.Background(), workspace.ID, ownerID)
	if err != nil {
		t.Fatalf("getting original owner membership: %v", err)
	}
	if err := u.UpdateMemberRole(context.Background(), ownerID, workspace.ID, originalOwner.ID, domain.RoleMember); err != nil {
		t.Fatalf("UpdateMemberRole() demote error = %v", err)
	}
}

func TestUpdateMemberRole_RefusesToLeaveOwnerless(t *testing.T) {
	u, _, members, _ := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	owner, err := members.GetByWorkspaceAndUser(context.Background(), workspace.ID, ownerID)
	if err != nil {
		t.Fatalf("getting owner membership: %v", err)
	}

	err = u.UpdateMemberRole(context.Background(), ownerID, workspace.ID, owner.ID, domain.RoleMember)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeConflict {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeConflict)
	}
}

func TestUpdateMemberRole_InvalidRole(t *testing.T) {
	u, _, members, _ := newTestWorkspaceUsecase()
	ownerID := uuid.New()

	workspace, err := u.Create(context.Background(), ownerID, "Acme Inc.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	owner, err := members.GetByWorkspaceAndUser(context.Background(), workspace.ID, ownerID)
	if err != nil {
		t.Fatalf("getting owner membership: %v", err)
	}

	err = u.UpdateMemberRole(context.Background(), ownerID, workspace.ID, owner.ID, domain.RoleAdmin)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
}
