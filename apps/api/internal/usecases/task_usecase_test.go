package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// mockTaskRepo is a hand-written in-memory domain.TaskRepository, per
// CONTRIBUTING.md's "mock repositories via interfaces" convention.
type mockTaskRepo struct {
	byID map[uuid.UUID]*domain.Task
}

func newMockTaskRepo() *mockTaskRepo {
	return &mockTaskRepo{byID: map[uuid.UUID]*domain.Task{}}
}

func (m *mockTaskRepo) Create(_ context.Context, task *domain.Task) error {
	task.ID = uuid.New()
	if task.Status == "" {
		task.Status = domain.StatusTodo
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt
	cp := *task
	m.byID[task.ID] = &cp
	return nil
}

func (m *mockTaskRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Task, error) {
	t, ok := m.byID[id]
	if !ok || t.DeletedAt != nil {
		return nil, domain.ErrTaskNotFound
	}
	cp := *t
	return &cp, nil
}

func (m *mockTaskRepo) ListByWorkspaceID(_ context.Context, workspaceID uuid.UUID, filter domain.TaskFilter) ([]*domain.Task, error) {
	var result []*domain.Task
	for _, t := range m.byID {
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
		cp := *t
		result = append(result, &cp)
	}
	return result, nil
}

func (m *mockTaskRepo) Update(_ context.Context, task *domain.Task) error {
	existing, ok := m.byID[task.ID]
	if !ok || existing.DeletedAt != nil {
		return domain.ErrTaskNotFound
	}
	cp := *task
	cp.UpdatedAt = time.Now()
	m.byID[task.ID] = &cp
	return nil
}

func (m *mockTaskRepo) Delete(_ context.Context, id uuid.UUID) error {
	t, ok := m.byID[id]
	if !ok || t.DeletedAt != nil {
		return domain.ErrTaskNotFound
	}
	now := time.Now()
	t.DeletedAt = &now
	return nil
}

func newTestTaskUsecase() (*TaskUsecase, *mockTaskRepo, *mockProjectRepo, *mockWorkspaceMemberRepo) {
	tasks := newMockTaskRepo()
	members := newMockWorkspaceMemberRepo()
	projects := newMockProjectRepo()
	u := NewTaskUsecase(tasks, projects, members)
	return u, tasks, projects, members
}

// seedProject inserts a project directly into projects (bypassing
// ProjectUsecase, which TaskUsecase doesn't depend on) so tests have a
// ProjectID to point tasks at.
func seedProject(t *testing.T, projects *mockProjectRepo, workspaceID uuid.UUID) *domain.Project {
	t.Helper()
	project := &domain.Project{WorkspaceID: workspaceID, Name: "Launch"}
	if err := projects.Create(context.Background(), project); err != nil {
		t.Fatalf("seeding project: %v", err)
	}
	return project
}

func TestTaskCreate_Success(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	project := seedProject(t, projects, workspaceID)

	task, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{
		ProjectID: project.ID,
		Title:     "  Ship it  ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if task.Title != "Ship it" {
		t.Errorf("Title = %q, want trimmed %q", task.Title, "Ship it")
	}
	if task.Status != domain.StatusTodo {
		t.Errorf("Status = %v, want %v", task.Status, domain.StatusTodo)
	}
	if task.WorkspaceID != workspaceID {
		t.Errorf("WorkspaceID = %v, want %v", task.WorkspaceID, workspaceID)
	}
}

func TestTaskCreate_EmptyTitle(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	project := seedProject(t, projects, workspaceID)

	_, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{ProjectID: project.ID, Title: "   "})
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
}

func TestTaskCreate_ForbiddenForNonMember(t *testing.T) {
	u, _, projects, _ := newTestTaskUsecase()
	workspaceID := uuid.New()
	project := seedProject(t, projects, workspaceID)

	_, err := u.Create(context.Background(), uuid.New(), workspaceID, CreateTaskInput{ProjectID: project.ID, Title: "Ship it"})
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeForbidden {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeForbidden)
	}
}

func TestTaskCreate_ProjectNotInWorkspace(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	seedMembership(t, members, workspaceA, callerID)
	projectInB := seedProject(t, projects, workspaceB)

	_, err := u.Create(context.Background(), callerID, workspaceA, CreateTaskInput{ProjectID: projectInB.ID, Title: "Ship it"})
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
	if appErr.Field != "projectId" {
		t.Errorf("Field = %q, want %q", appErr.Field, "projectId")
	}
}

func TestTaskCreate_AssigneeNotMember(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceID, callerID, strangerID := uuid.New(), uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	project := seedProject(t, projects, workspaceID)

	_, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{
		ProjectID:  project.ID,
		Title:      "Ship it",
		AssigneeID: &strangerID,
	})
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
	if appErr.Field != "assigneeId" {
		t.Errorf("Field = %q, want %q", appErr.Field, "assigneeId")
	}
}

func TestTaskGetByID_CrossWorkspaceIsNotFound(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	seedMembership(t, members, workspaceA, callerID)
	seedMembership(t, members, workspaceB, callerID)
	project := seedProject(t, projects, workspaceB)

	task, err := u.Create(context.Background(), callerID, workspaceB, CreateTaskInput{ProjectID: project.ID, Title: "Ship it"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = u.GetByID(context.Background(), callerID, workspaceA, task.ID)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v — a task must not be reachable through a different workspace's ID", appErr.Code, apperrors.CodeNotFound)
	}
}

func TestTaskListByWorkspaceID_Filters(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceID, callerID, assigneeID := uuid.New(), uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	seedMembership(t, members, workspaceID, assigneeID)
	projectA := seedProject(t, projects, workspaceID)
	projectB := seedProject(t, projects, workspaceID)

	taskA, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{ProjectID: projectA.ID, Title: "In project A", AssigneeID: &assigneeID})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{ProjectID: projectB.ID, Title: "In project B"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	byProject, err := u.ListByWorkspaceID(context.Background(), callerID, workspaceID, domain.TaskFilter{ProjectID: &projectA.ID})
	if err != nil {
		t.Fatalf("ListByWorkspaceID() error = %v", err)
	}
	if len(byProject) != 1 || byProject[0].ID != taskA.ID {
		t.Errorf("filtering by projectId: got %d tasks, want just taskA", len(byProject))
	}

	byAssignee, err := u.ListByWorkspaceID(context.Background(), callerID, workspaceID, domain.TaskFilter{AssigneeID: &assigneeID})
	if err != nil {
		t.Fatalf("ListByWorkspaceID() error = %v", err)
	}
	if len(byAssignee) != 1 || byAssignee[0].ID != taskA.ID {
		t.Errorf("filtering by assigneeId: got %d tasks, want just taskA", len(byAssignee))
	}

	doneStatus := domain.StatusDone
	byStatus, err := u.ListByWorkspaceID(context.Background(), callerID, workspaceID, domain.TaskFilter{Status: &doneStatus})
	if err != nil {
		t.Fatalf("ListByWorkspaceID() error = %v", err)
	}
	if len(byStatus) != 0 {
		t.Errorf("filtering by status=done: got %d tasks, want 0 (nothing is done yet)", len(byStatus))
	}
}

func TestTaskUpdate_Success(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	project := seedProject(t, projects, workspaceID)

	task, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{ProjectID: project.ID, Title: "Ship it"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	newTitle := "Ship it faster"
	inProgress := domain.StatusInProgress
	updated, err := u.Update(context.Background(), callerID, workspaceID, task.ID, UpdateTaskInput{
		Title:  &newTitle,
		Status: &inProgress,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Title != newTitle {
		t.Errorf("Title = %q, want %q", updated.Title, newTitle)
	}
	if updated.Status != domain.StatusInProgress {
		t.Errorf("Status = %v, want %v", updated.Status, domain.StatusInProgress)
	}
}

func TestTaskUpdate_ProjectChangeValidatesNewProject(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	seedMembership(t, members, workspaceA, callerID)
	projectInA := seedProject(t, projects, workspaceA)
	projectInB := seedProject(t, projects, workspaceB)

	task, err := u.Create(context.Background(), callerID, workspaceA, CreateTaskInput{ProjectID: projectInA.ID, Title: "Ship it"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = u.Update(context.Background(), callerID, workspaceA, task.ID, UpdateTaskInput{ProjectID: &projectInB.ID})
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
}

func TestTaskUpdate_ClearAssignee(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	project := seedProject(t, projects, workspaceID)

	task, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{
		ProjectID:  project.ID,
		Title:      "Ship it",
		AssigneeID: &callerID,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	updated, err := u.Update(context.Background(), callerID, workspaceID, task.ID, UpdateTaskInput{ClearAssignee: true})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.AssigneeID != nil {
		t.Errorf("AssigneeID = %v, want nil after ClearAssignee", updated.AssigneeID)
	}
}

func TestTaskUpdate_InvalidStatus(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	project := seedProject(t, projects, workspaceID)

	task, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{ProjectID: project.ID, Title: "Ship it"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	bogus := domain.TaskStatus("archived")
	_, err = u.Update(context.Background(), callerID, workspaceID, task.ID, UpdateTaskInput{Status: &bogus})
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeValidation {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
	}
}

func TestTaskDelete_Success(t *testing.T) {
	u, _, projects, members := newTestTaskUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)
	project := seedProject(t, projects, workspaceID)

	task, err := u.Create(context.Background(), callerID, workspaceID, CreateTaskInput{ProjectID: project.ID, Title: "Ship it"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := u.Delete(context.Background(), callerID, workspaceID, task.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = u.GetByID(context.Background(), callerID, workspaceID, task.ID)
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v — a deleted task should read as not found", appErr.Code, apperrors.CodeNotFound)
	}
}

func TestTaskDelete_NotFound(t *testing.T) {
	u, _, _, members := newTestTaskUsecase()
	workspaceID, callerID := uuid.New(), uuid.New()
	seedMembership(t, members, workspaceID, callerID)

	err := u.Delete(context.Background(), callerID, workspaceID, uuid.New())
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeNotFound {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeNotFound)
	}
}
