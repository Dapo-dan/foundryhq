package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/middleware"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
)

// newTestJWTManager and bearerToken let workspace-scoped handler tests run
// requests through the real middleware.Auth (a generated, genuinely valid
// access token) rather than reaching into its unexported context key — the
// same authentication path a real request takes.
func newTestJWTManager() *jwt.Manager {
	return jwt.NewManager("access-secret", "refresh-secret", 15*time.Minute, 168*time.Hour)
}

func bearerToken(t *testing.T, manager *jwt.Manager, userID uuid.UUID) string {
	t.Helper()
	token, err := manager.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("generating access token: %v", err)
	}
	return token
}

func newProtectedTestRouter(manager *jwt.Manager, register func(group *gin.RouterGroup)) *gin.Engine {
	router := gin.New()
	register(router.Group("", middleware.Auth(manager)))
	return router
}

func doAuthedJSONRequest(router *gin.Engine, method, path, token string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// jsonRequestNoAuth builds a request/recorder pair with no Authorization
// header, for asserting middleware.Auth itself rejects an anonymous call —
// the caller still has to invoke router.ServeHTTP(w, req).
func jsonRequestNoAuth(method, path string, body any) (*http.Request, *httptest.ResponseRecorder) {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

// assertErrorCode fails t unless w's body is a { "error": { "code": ... } }
// envelope matching want.
func assertErrorCode(t *testing.T, w *httptest.ResponseRecorder, want apperrors.Code) {
	t.Helper()
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v, body = %s", err, w.Body.String())
	}
	if body.Error.Code != string(want) {
		t.Errorf("error.code = %q, want %q", body.Error.Code, want)
	}
}

// This file collects the in-memory domain.*Repository fakes shared by
// workspace_handler_test.go, project_handler_test.go, task_handler_test.go,
// and sprint_handler_test.go — every one of those resources is
// workspace-scoped, so they all need a WorkspaceMemberRepository fake at
// minimum. Kept in one place rather than redefined per file, mirroring
// auth_handler_test.go's "local fakes, not shared usecase test doubles"
// convention, just shared across this package's own handler tests instead
// of duplicated four times.

type fakeWorkspaceRepo struct {
	byID   map[uuid.UUID]*domain.Workspace
	bySlug map[string]bool
}

func newFakeWorkspaceRepo() *fakeWorkspaceRepo {
	return &fakeWorkspaceRepo{byID: map[uuid.UUID]*domain.Workspace{}, bySlug: map[string]bool{}}
}

func (r *fakeWorkspaceRepo) Create(_ context.Context, workspace *domain.Workspace, owner *domain.WorkspaceMember) error {
	workspace.ID = uuid.New()
	workspace.CreatedAt = time.Now()
	workspace.UpdatedAt = workspace.CreatedAt
	r.byID[workspace.ID] = workspace
	r.bySlug[workspace.Slug] = true
	owner.ID = uuid.New()
	owner.WorkspaceID = workspace.ID
	owner.InvitedAt = time.Now()
	return nil
}

func (r *fakeWorkspaceRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Workspace, error) {
	w, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrWorkspaceNotFound
	}
	return w, nil
}

func (r *fakeWorkspaceRepo) Update(_ context.Context, workspace *domain.Workspace) error {
	if _, ok := r.byID[workspace.ID]; !ok {
		return domain.ErrWorkspaceNotFound
	}
	workspace.UpdatedAt = time.Now()
	r.byID[workspace.ID] = workspace
	r.bySlug[workspace.Slug] = true
	return nil
}

func (r *fakeWorkspaceRepo) SlugExists(_ context.Context, slug string) (bool, error) {
	return r.bySlug[slug], nil
}

func (r *fakeWorkspaceRepo) ListForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Workspace, error) {
	return nil, nil // not exercised by handler tests — WorkspaceUsecase.ListForUser doesn't need membership fan-out here.
}

// fakeWorkspaceMemberRepo backs every resource's membership guard
// (requireMembership/requireOwner). seedMember is the test-facing way to
// grant a user a role in a workspace before exercising a handler.
type fakeWorkspaceMemberRepo struct {
	byID map[uuid.UUID]*domain.WorkspaceMember
}

func newFakeWorkspaceMemberRepo() *fakeWorkspaceMemberRepo {
	return &fakeWorkspaceMemberRepo{byID: map[uuid.UUID]*domain.WorkspaceMember{}}
}

func (r *fakeWorkspaceMemberRepo) seedMember(workspaceID, userID uuid.UUID, role domain.WorkspaceRole) *domain.WorkspaceMember {
	now := time.Now()
	member := &domain.WorkspaceMember{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		InvitedAt:   now,
		JoinedAt:    &now,
	}
	r.byID[member.ID] = member
	return member
}

func (r *fakeWorkspaceMemberRepo) Create(_ context.Context, member *domain.WorkspaceMember) error {
	member.ID = uuid.New()
	member.InvitedAt = time.Now()
	r.byID[member.ID] = member
	return nil
}

func (r *fakeWorkspaceMemberRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.WorkspaceMember, error) {
	m, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrWorkspaceMemberNotFound
	}
	return m, nil
}

func (r *fakeWorkspaceMemberRepo) MarkJoined(_ context.Context, id uuid.UUID) error {
	m, ok := r.byID[id]
	if !ok {
		return domain.ErrWorkspaceMemberNotFound
	}
	now := time.Now()
	m.JoinedAt = &now
	return nil
}

func (r *fakeWorkspaceMemberRepo) GetByWorkspaceAndUser(_ context.Context, workspaceID, userID uuid.UUID) (*domain.WorkspaceMember, error) {
	for _, m := range r.byID {
		if m.WorkspaceID == workspaceID && m.UserID == userID {
			return m, nil
		}
	}
	return nil, domain.ErrWorkspaceMemberNotFound
}

func (r *fakeWorkspaceMemberRepo) ListByWorkspaceIDWithUser(_ context.Context, workspaceID uuid.UUID) ([]*domain.WorkspaceMemberWithUser, error) {
	var result []*domain.WorkspaceMemberWithUser
	for _, m := range r.byID {
		if m.WorkspaceID == workspaceID {
			result = append(result, &domain.WorkspaceMemberWithUser{WorkspaceMember: *m})
		}
	}
	return result, nil
}

func (r *fakeWorkspaceMemberRepo) UpdateRole(_ context.Context, id uuid.UUID, role domain.WorkspaceRole) error {
	m, ok := r.byID[id]
	if !ok {
		return domain.ErrWorkspaceMemberNotFound
	}
	m.Role = role
	return nil
}

type fakeProjectRepo struct {
	byID map[uuid.UUID]*domain.Project
}

func newFakeProjectRepo() *fakeProjectRepo {
	return &fakeProjectRepo{byID: map[uuid.UUID]*domain.Project{}}
}

func (r *fakeProjectRepo) seedProject(workspaceID uuid.UUID, name string) *domain.Project {
	now := time.Now()
	project := &domain.Project{ID: uuid.New(), WorkspaceID: workspaceID, Name: name, CreatedAt: now, UpdatedAt: now}
	r.byID[project.ID] = project
	return project
}

func (r *fakeProjectRepo) Create(_ context.Context, project *domain.Project) error {
	project.ID = uuid.New()
	project.CreatedAt = time.Now()
	project.UpdatedAt = project.CreatedAt
	r.byID[project.ID] = project
	return nil
}

func (r *fakeProjectRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Project, error) {
	p, ok := r.byID[id]
	if !ok || p.DeletedAt != nil {
		return nil, domain.ErrProjectNotFound
	}
	return p, nil
}

func (r *fakeProjectRepo) ListByWorkspaceID(_ context.Context, workspaceID uuid.UUID) ([]*domain.Project, error) {
	var result []*domain.Project
	for _, p := range r.byID {
		if p.WorkspaceID == workspaceID && p.DeletedAt == nil {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *fakeProjectRepo) Update(_ context.Context, project *domain.Project) error {
	if _, ok := r.byID[project.ID]; !ok {
		return domain.ErrProjectNotFound
	}
	project.UpdatedAt = time.Now()
	r.byID[project.ID] = project
	return nil
}

func (r *fakeProjectRepo) Delete(_ context.Context, id uuid.UUID) error {
	p, ok := r.byID[id]
	if !ok || p.DeletedAt != nil {
		return domain.ErrProjectNotFound
	}
	now := time.Now()
	p.DeletedAt = &now
	return nil
}

type fakeTaskRepo struct {
	byID map[uuid.UUID]*domain.Task
}

func newFakeTaskRepo() *fakeTaskRepo {
	return &fakeTaskRepo{byID: map[uuid.UUID]*domain.Task{}}
}

func (r *fakeTaskRepo) Create(_ context.Context, task *domain.Task) error {
	task.ID = uuid.New()
	if task.Status == "" {
		task.Status = domain.StatusTodo
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt
	r.byID[task.ID] = task
	return nil
}

func (r *fakeTaskRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Task, error) {
	t, ok := r.byID[id]
	if !ok || t.DeletedAt != nil {
		return nil, domain.ErrTaskNotFound
	}
	return t, nil
}

func (r *fakeTaskRepo) ListByWorkspaceID(_ context.Context, workspaceID uuid.UUID, filter domain.TaskFilter) ([]*domain.Task, error) {
	var result []*domain.Task
	for _, t := range r.byID {
		if t.WorkspaceID != workspaceID || t.DeletedAt != nil {
			continue
		}
		if filter.ProjectID != nil && t.ProjectID != *filter.ProjectID {
			continue
		}
		if filter.Status != nil && t.Status != *filter.Status {
			continue
		}
		if filter.AssigneeID != nil && (t.AssigneeID == nil || *t.AssigneeID != *filter.AssigneeID) {
			continue
		}
		if filter.SprintID != nil && (t.SprintID == nil || *t.SprintID != *filter.SprintID) {
			continue
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *fakeTaskRepo) Update(_ context.Context, task *domain.Task) error {
	existing, ok := r.byID[task.ID]
	if !ok || existing.DeletedAt != nil {
		return domain.ErrTaskNotFound
	}
	task.UpdatedAt = time.Now()
	r.byID[task.ID] = task
	return nil
}

func (r *fakeTaskRepo) Delete(_ context.Context, id uuid.UUID) error {
	t, ok := r.byID[id]
	if !ok || t.DeletedAt != nil {
		return domain.ErrTaskNotFound
	}
	now := time.Now()
	t.DeletedAt = &now
	return nil
}

func (r *fakeTaskRepo) SumStoryPointsForSprint(_ context.Context, sprintID uuid.UUID, startDate, endDate time.Time) (int, error) {
	sum := 0
	exclusiveEnd := endDate.AddDate(0, 0, 1)
	for _, t := range r.byID {
		if t.DeletedAt != nil || t.SprintID == nil || *t.SprintID != sprintID || t.Status != domain.StatusDone {
			continue
		}
		if t.UpdatedAt.Before(startDate) || !t.UpdatedAt.Before(exclusiveEnd) {
			continue
		}
		if t.StoryPoints != nil {
			sum += *t.StoryPoints
		}
	}
	return sum, nil
}

// fakeInviteTokenRepo is a hand-written in-memory domain.InviteTokenRepository.
type fakeInviteTokenRepo struct {
	byHash map[string]*domain.InviteToken
}

func newFakeInviteTokenRepo() *fakeInviteTokenRepo {
	return &fakeInviteTokenRepo{byHash: map[string]*domain.InviteToken{}}
}

func (r *fakeInviteTokenRepo) Create(_ context.Context, token *domain.InviteToken) error {
	token.ID = uuid.New()
	token.CreatedAt = time.Now()
	r.byHash[token.TokenHash] = token
	return nil
}

func (r *fakeInviteTokenRepo) GetByTokenHash(_ context.Context, tokenHash string) (*domain.InviteToken, error) {
	t, ok := r.byHash[tokenHash]
	if !ok {
		return nil, domain.ErrInviteTokenNotFound
	}
	return t, nil
}

func (r *fakeInviteTokenRepo) MarkUsed(_ context.Context, tokenHash string) error {
	t, ok := r.byHash[tokenHash]
	if !ok {
		return nil
	}
	now := time.Now()
	t.UsedAt = &now
	return nil
}

type fakeSprintRepo struct {
	byID map[uuid.UUID]*domain.Sprint
}

func newFakeSprintRepo() *fakeSprintRepo {
	return &fakeSprintRepo{byID: map[uuid.UUID]*domain.Sprint{}}
}

func (r *fakeSprintRepo) seedSprint(workspaceID uuid.UUID, name string, start, end time.Time) *domain.Sprint {
	now := time.Now()
	sprint := &domain.Sprint{ID: uuid.New(), WorkspaceID: workspaceID, Name: name, StartDate: start, EndDate: end, CreatedAt: now, UpdatedAt: now}
	r.byID[sprint.ID] = sprint
	return sprint
}

func (r *fakeSprintRepo) Create(_ context.Context, sprint *domain.Sprint) error {
	sprint.ID = uuid.New()
	sprint.CreatedAt = time.Now()
	sprint.UpdatedAt = sprint.CreatedAt
	r.byID[sprint.ID] = sprint
	return nil
}

func (r *fakeSprintRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Sprint, error) {
	s, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrSprintNotFound
	}
	return s, nil
}

func (r *fakeSprintRepo) ListByWorkspaceID(_ context.Context, workspaceID uuid.UUID) ([]*domain.Sprint, error) {
	var result []*domain.Sprint
	for _, s := range r.byID {
		if s.WorkspaceID == workspaceID {
			result = append(result, s)
		}
	}
	return result, nil
}
