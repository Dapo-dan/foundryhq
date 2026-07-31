package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// mockProjectRepo is a hand-written in-memory domain.ProjectRepository, per
// CONTRIBUTING.md's "mock repositories via interfaces" convention.
type mockProjectRepo struct {
	byID map[uuid.UUID]*domain.Project
}

func newMockProjectRepo() *mockProjectRepo {
	return &mockProjectRepo{byID: map[uuid.UUID]*domain.Project{}}
}

func (m *mockProjectRepo) Create(_ context.Context, project *domain.Project) error {
	project.ID = uuid.New()
	project.CreatedAt = time.Now()
	project.UpdatedAt = project.CreatedAt
	cp := *project
	m.byID[project.ID] = &cp
	return nil
}

func (m *mockProjectRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Project, error) {
	p, ok := m.byID[id]
	if !ok || p.DeletedAt != nil {
		return nil, domain.ErrProjectNotFound
	}
	cp := *p
	return &cp, nil
}

func (m *mockProjectRepo) ListByWorkspaceID(_ context.Context, workspaceID uuid.UUID) ([]*domain.Project, error) {
	var result []*domain.Project
	for _, p := range m.byID {
		if p.WorkspaceID == workspaceID && p.DeletedAt == nil {
			cp := *p
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (m *mockProjectRepo) Update(_ context.Context, project *domain.Project) error {
	existing, ok := m.byID[project.ID]
	if !ok || existing.DeletedAt != nil {
		return domain.ErrProjectNotFound
	}
	cp := *project
	cp.UpdatedAt = time.Now()
	m.byID[project.ID] = &cp
	return nil
}

func (m *mockProjectRepo) Delete(_ context.Context, id uuid.UUID) error {
	p, ok := m.byID[id]
	if !ok || p.DeletedAt != nil {
		return domain.ErrProjectNotFound
	}
	now := time.Now()
	p.DeletedAt = &now
	return nil
}

func newTestProjectUsecase() (*ProjectUsecase, *mockProjectRepo, *mockWorkspaceMemberRepo) {
	projects := newMockProjectRepo()
	members := newMockWorkspaceMemberRepo()
	u := NewProjectUsecase(projects, members)
	return u, projects, members
}

// seedMembership makes userID a member of workspaceID so requireMembership
// passes for it.
func seedMembership(t *testing.T, members *mockWorkspaceMemberRepo, workspaceID, userID uuid.UUID) {
	t.Helper()
	if err := members.Create(context.Background(), &domain.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        domain.RoleMember,
	}); err != nil {
		t.Fatalf("seeding membership: %v", err)
	}
}

func TestProjectCreate_Success(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	description := "Q3 roadmap"
	project, err := u.Create(context.Background(), callerID, workspaceID, "  Launch  ", &description)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if project.Name != "Launch" {
		t.Errorf("Name = %q, want trimmed %q", project.Name, "Launch")
	}
	if project.Description == nil || *project.Description != description {
		t.Errorf("Description = %v, want %q", project.Description, description)
	}
	if project.WorkspaceID != workspaceID {
		t.Errorf("WorkspaceID = %v, want %v", project.WorkspaceID, workspaceID)
	}
}

func TestProjectCreate_EmptyName(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	_, err := u.Create(context.Background(), callerID, workspaceID, "   ", nil)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
}

func TestProjectCreate_ForbiddenForNonMember(t *testing.T) {
	u, _, _ := newTestProjectUsecase()

	_, err := u.Create(context.Background(), uuid.New(), uuid.New(), "Launch", nil)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeForbidden {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeForbidden)
	}
}

func TestProjectGetByID_NotFound(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	_, err := u.GetByID(context.Background(), callerID, workspaceID, uuid.New())
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeNotFound)
	}
}

func TestProjectGetByID_CrossWorkspaceIsNotFound(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	// Caller is a member of both workspaces, but the project belongs to B.
	seedMembership(t, members, workspaceA, callerID)
	seedMembership(t, members, workspaceB, callerID)

	project, err := u.Create(context.Background(), callerID, workspaceB, "Launch", nil)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = u.GetByID(context.Background(), callerID, workspaceA, project.ID)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v — a project must not be reachable through a different workspace's ID", appErr.Code, apperrors.CodeNotFound)
	}
}

func TestProjectListByWorkspaceID(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceID, otherWorkspaceID, callerID := uuid.New(), uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	seedMembership(t, members, otherWorkspaceID, callerID)

	if _, err := u.Create(context.Background(), callerID, workspaceID, "Launch", nil); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := u.Create(context.Background(), callerID, otherWorkspaceID, "Other", nil); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	projects, err := u.ListByWorkspaceID(context.Background(), callerID, workspaceID)
	if err != nil {
		t.Fatalf("ListByWorkspaceID() error = %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Name != "Launch" {
		t.Errorf("Name = %q, want %q", projects[0].Name, "Launch")
	}
}

func TestProjectUpdate_Success(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	project, err := u.Create(context.Background(), callerID, workspaceID, "Launch", nil)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	newName := "Launch v2"
	newDescription := "Updated scope"
	updated, err := u.Update(context.Background(), callerID, workspaceID, project.ID, UpdateProjectInput{
		Name:        &newName,
		Description: &newDescription,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Name != newName {
		t.Errorf("Name = %q, want %q", updated.Name, newName)
	}
	if updated.Description == nil || *updated.Description != newDescription {
		t.Errorf("Description = %v, want %q", updated.Description, newDescription)
	}
}

func TestProjectUpdate_ClearDescription(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	description := "Initial"
	project, err := u.Create(context.Background(), callerID, workspaceID, "Launch", &description)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	blank := "   "
	updated, err := u.Update(context.Background(), callerID, workspaceID, project.ID, UpdateProjectInput{Description: &blank})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Description != nil {
		t.Errorf("Description = %v, want nil after clearing with a blank string", updated.Description)
	}
}

func TestProjectDelete_Success(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	project, err := u.Create(context.Background(), callerID, workspaceID, "Launch", nil)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := u.Delete(context.Background(), callerID, workspaceID, project.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = u.GetByID(context.Background(), callerID, workspaceID, project.ID)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v — a deleted project should read as not found", appErr.Code, apperrors.CodeNotFound)
	}

	projects, err := u.ListByWorkspaceID(context.Background(), callerID, workspaceID)
	if err != nil {
		t.Fatalf("ListByWorkspaceID() error = %v", err)
	}
	if len(projects) != 0 {
		t.Errorf("expected a deleted project to be excluded from the list, got %d", len(projects))
	}
}

func TestProjectDelete_NotFound(t *testing.T) {
	u, _, members := newTestProjectUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	err := u.Delete(context.Background(), callerID, workspaceID, uuid.New())
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeNotFound)
	}
}
