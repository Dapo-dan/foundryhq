package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// mockSprintRepo is a hand-written in-memory domain.SprintRepository, per
// CONTRIBUTING.md's "mock repositories via interfaces" convention.
type mockSprintRepo struct {
	byID map[uuid.UUID]*domain.Sprint
}

func newMockSprintRepo() *mockSprintRepo {
	return &mockSprintRepo{byID: map[uuid.UUID]*domain.Sprint{}}
}

func (m *mockSprintRepo) Create(_ context.Context, sprint *domain.Sprint) error {
	sprint.ID = uuid.New()
	sprint.CreatedAt = time.Now()
	sprint.UpdatedAt = sprint.CreatedAt
	cp := *sprint
	m.byID[sprint.ID] = &cp
	return nil
}

func (m *mockSprintRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Sprint, error) {
	s, ok := m.byID[id]
	if !ok {
		return nil, domain.ErrSprintNotFound
	}
	cp := *s
	return &cp, nil
}

func (m *mockSprintRepo) ListByWorkspaceID(_ context.Context, workspaceID uuid.UUID) ([]*domain.Sprint, error) {
	var result []*domain.Sprint
	for _, s := range m.byID {
		if s.WorkspaceID == workspaceID {
			cp := *s
			result = append(result, &cp)
		}
	}
	return result, nil
}

func newTestSprintUsecase() (*SprintUsecase, *mockSprintRepo, *mockTaskRepo, *mockWorkspaceMemberRepo) {
	sprints := newMockSprintRepo()
	tasks := newMockTaskRepo()
	members := newMockWorkspaceMemberRepo()
	u := NewSprintUsecase(sprints, tasks, members)
	return u, sprints, tasks, members
}

func intPtr(n int) *int { return &n }

func TestSprintCreate_Success(t *testing.T) {
	u, _, _, members := newTestSprintUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprint, err := u.Create(context.Background(), callerID, workspaceID, "  Sprint 1  ", start, end)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if sprint.Name != "Sprint 1" {
		t.Errorf("Name = %q, want trimmed %q", sprint.Name, "Sprint 1")
	}
	if sprint.WorkspaceID != workspaceID {
		t.Errorf("WorkspaceID = %v, want %v", sprint.WorkspaceID, workspaceID)
	}
}

func TestSprintCreate_EmptyName(t *testing.T) {
	u, _, _, members := newTestSprintUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	_, err := u.Create(context.Background(), callerID, workspaceID, "   ", start, start)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
}

func TestSprintCreate_EndBeforeStart(t *testing.T) {
	u, _, _, members := newTestSprintUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	start := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	_, err := u.Create(context.Background(), callerID, workspaceID, "Sprint 1", start, end)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
	if appErr.Field != "endDate" {
		t.Errorf("Field = %q, want %q", appErr.Field, "endDate")
	}
}

func TestSprintCreate_ForbiddenForNonMember(t *testing.T) {
	u, _, _, _ := newTestSprintUsecase()
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)

	_, err := u.Create(context.Background(), uuid.New(), uuid.New(), "Sprint 1", start, end)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeForbidden {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeForbidden)
	}
}

func TestSprintGetByID_WithTasks(t *testing.T) {
	u, _, tasks, members := newTestSprintUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprint, err := u.Create(context.Background(), callerID, workspaceID, "Sprint 1", start, end)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	inSprint := &domain.Task{ID: uuid.New(), WorkspaceID: workspaceID, SprintID: &sprint.ID, Title: "In sprint", Status: domain.StatusTodo}
	tasks.byID[inSprint.ID] = inSprint
	outOfSprint := &domain.Task{ID: uuid.New(), WorkspaceID: workspaceID, Title: "Backlog", Status: domain.StatusTodo}
	tasks.byID[outOfSprint.ID] = outOfSprint

	result, err := u.GetByID(context.Background(), callerID, workspaceID, sprint.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if result.Sprint.ID != sprint.ID {
		t.Errorf("Sprint.ID = %v, want %v", result.Sprint.ID, sprint.ID)
	}
	if len(result.Tasks) != 1 || result.Tasks[0].ID != inSprint.ID {
		t.Errorf("expected exactly the in-sprint task, got %d tasks", len(result.Tasks))
	}
}

func TestSprintGetByID_CrossWorkspaceIsNotFound(t *testing.T) {
	u, _, _, members := newTestSprintUsecase()
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	seedMembership(t, members, workspaceA, callerID)
	seedMembership(t, members, workspaceB, callerID)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprint, err := u.Create(context.Background(), callerID, workspaceB, "Sprint 1", start, end)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = u.GetByID(context.Background(), callerID, workspaceA, sprint.ID)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v — a sprint must not be reachable through a different workspace's ID", appErr.Code, apperrors.CodeNotFound)
	}
}

func TestSprintListForWorkspace(t *testing.T) {
	u, _, _, members := newTestSprintUsecase()
	workspaceID, otherWorkspaceID, callerID := uuid.New(), uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	seedMembership(t, members, otherWorkspaceID, callerID)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	if _, err := u.Create(context.Background(), callerID, workspaceID, "Sprint 1", start, end); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := u.Create(context.Background(), callerID, otherWorkspaceID, "Other Sprint", start, end); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	sprints, err := u.ListForWorkspace(context.Background(), callerID, workspaceID)
	if err != nil {
		t.Fatalf("ListForWorkspace() error = %v", err)
	}
	if len(sprints) != 1 || sprints[0].Name != "Sprint 1" {
		t.Errorf("expected exactly 1 sprint named %q, got %d", "Sprint 1", len(sprints))
	}
}

// TestSprintGetVelocity_SumsOnlyDoneWithinRange is the core acceptance-
// criteria test (.ai/business-analysis/acceptance-criteria.md): velocity
// counts only done tasks whose updated_at falls within the sprint's date
// range: tasks done early, done late (after the sprint closes), or not
// done at all must all be excluded.
func TestSprintGetVelocity_SumsOnlyDoneWithinRange(t *testing.T) {
	u, _, tasks, members := newTestSprintUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprint, err := u.Create(context.Background(), callerID, workspaceID, "Sprint 1", start, end)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	withinRange := &domain.Task{
		ID: uuid.New(), WorkspaceID: workspaceID, SprintID: &sprint.ID,
		Status: domain.StatusDone, StoryPoints: intPtr(5),
		UpdatedAt: time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC),
	}
	onTheClosingDay := &domain.Task{
		ID: uuid.New(), WorkspaceID: workspaceID, SprintID: &sprint.ID,
		Status: domain.StatusDone, StoryPoints: intPtr(3),
		UpdatedAt: time.Date(2026, 8, 14, 23, 0, 0, 0, time.UTC), // late on end_date itself — must still count
	}
	doneBeforeSprint := &domain.Task{
		ID: uuid.New(), WorkspaceID: workspaceID, SprintID: &sprint.ID,
		Status: domain.StatusDone, StoryPoints: intPtr(8),
		UpdatedAt: time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
	}
	doneAfterSprintCloses := &domain.Task{
		ID: uuid.New(), WorkspaceID: workspaceID, SprintID: &sprint.ID,
		Status: domain.StatusDone, StoryPoints: intPtr(13),
		UpdatedAt: time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC),
	}
	notDone := &domain.Task{
		ID: uuid.New(), WorkspaceID: workspaceID, SprintID: &sprint.ID,
		Status: domain.StatusInProgress, StoryPoints: intPtr(21),
		UpdatedAt: time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC),
	}
	for _, task := range []*domain.Task{withinRange, onTheClosingDay, doneBeforeSprint, doneAfterSprintCloses, notDone} {
		tasks.byID[task.ID] = task
	}

	velocity, err := u.GetVelocity(context.Background(), callerID, workspaceID, sprint.ID)
	if err != nil {
		t.Fatalf("GetVelocity() error = %v", err)
	}
	if velocity != 8 {
		t.Errorf("GetVelocity() = %d, want 8 (5 + 3, excluding before/after/not-done)", velocity)
	}
}

func TestSprintGetVelocity_ZeroWhenNothingDone(t *testing.T) {
	u, _, _, members := newTestSprintUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprint, err := u.Create(context.Background(), callerID, workspaceID, "Sprint 1", start, end)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	velocity, err := u.GetVelocity(context.Background(), callerID, workspaceID, sprint.ID)
	if err != nil {
		t.Fatalf("GetVelocity() error = %v", err)
	}
	if velocity != 0 {
		t.Errorf("GetVelocity() = %d, want 0", velocity)
	}
}
