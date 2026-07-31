package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
)

func newTestTaskHandler() (*TaskHandler, *fakeTaskRepo, *fakeProjectRepo, *fakeWorkspaceMemberRepo, *jwt.Manager) {
	tasks := newFakeTaskRepo()
	projects := newFakeProjectRepo()
	members := newFakeWorkspaceMemberRepo()
	sprints := newFakeSprintRepo()
	usecase := usecases.NewTaskUsecase(tasks, projects, sprints, members)
	return NewTaskHandler(usecase), tasks, projects, members, newTestJWTManager()
}

func newTaskTestRouter(manager *jwt.Manager, h *TaskHandler) *gin.Engine {
	return newProtectedTestRouter(manager, func(protected *gin.RouterGroup) {
		group := protected.Group("/workspaces/:workspaceId/tasks")
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:taskId", h.Get)
		group.PATCH("/:taskId", h.Update)
		group.DELETE("/:taskId", h.Delete)
	})
}

func TestTaskHandler_Create_Success(t *testing.T) {
	h, _, projects, members, manager := newTestTaskHandler()
	router := newTaskTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)
	project := projects.seedProject(workspaceID, "Launch")

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces/"+workspaceID.String()+"/tasks", bearerToken(t, manager, callerID), map[string]any{
		"projectId": project.ID.String(),
		"title":     "Ship it",
	})

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusCreated, w.Body.String())
	}
	var body struct {
		Data taskResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Data.Priority != "medium" {
		t.Errorf("priority = %q, want default %q", body.Data.Priority, "medium")
	}
	if body.Data.Status != "todo" {
		t.Errorf("status = %q, want default %q", body.Data.Status, "todo")
	}
}

func TestTaskHandler_Create_ProjectNotInWorkspace(t *testing.T) {
	h, _, projects, members, manager := newTestTaskHandler()
	router := newTaskTestRouter(manager, h)
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	members.seedMember(workspaceA, callerID, domain.RoleMember)
	projectInB := projects.seedProject(workspaceB, "Other workspace's project")

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces/"+workspaceA.String()+"/tasks", bearerToken(t, manager, callerID), map[string]any{
		"projectId": projectInB.ID.String(),
		"title":     "Ship it",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
	assertErrorCode(t, w, apperrors.CodeValidation)
}

func TestTaskHandler_Get_CrossWorkspaceIsNotFound(t *testing.T) {
	h, tasks, projects, members, manager := newTestTaskHandler()
	router := newTaskTestRouter(manager, h)
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	members.seedMember(workspaceA, callerID, domain.RoleMember)
	members.seedMember(workspaceB, callerID, domain.RoleMember)
	project := projects.seedProject(workspaceB, "Launch")
	task := &domain.Task{WorkspaceID: workspaceB, ProjectID: project.ID, Title: "Ship it"}
	if err := tasks.Create(t.Context(), task); err != nil {
		t.Fatalf("seeding task: %v", err)
	}

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceA.String()+"/tasks/"+task.ID.String(), bearerToken(t, manager, callerID), nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body = %s — a task must not be reachable through a different workspace's path", w.Code, http.StatusNotFound, w.Body.String())
	}
}

func TestTaskHandler_Update_ClearAssignee(t *testing.T) {
	h, tasks, projects, members, manager := newTestTaskHandler()
	router := newTaskTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)
	project := projects.seedProject(workspaceID, "Launch")
	task := &domain.Task{WorkspaceID: workspaceID, ProjectID: project.ID, Title: "Ship it", AssigneeID: &callerID}
	if err := tasks.Create(t.Context(), task); err != nil {
		t.Fatalf("seeding task: %v", err)
	}

	w := doAuthedJSONRequest(router, http.MethodPatch, "/workspaces/"+workspaceID.String()+"/tasks/"+task.ID.String(), bearerToken(t, manager, callerID), map[string]any{
		"clearAssignee": true,
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
	var body struct {
		Data taskResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Data.AssigneeID != nil {
		t.Errorf("assigneeId = %v, want nil after clearAssignee", body.Data.AssigneeID)
	}
}

func TestTaskHandler_Update_ClearSprintAndDueDate(t *testing.T) {
	h, tasks, projects, members, manager := newTestTaskHandler()
	router := newTaskTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)
	project := projects.seedProject(workspaceID, "Launch")

	dueDate, err := time.Parse(dateOnlyLayout, "2026-08-15")
	if err != nil {
		t.Fatalf("parsing seed due date: %v", err)
	}
	sprintID := uuid.New()
	task := &domain.Task{WorkspaceID: workspaceID, ProjectID: project.ID, Title: "Ship it", SprintID: &sprintID, DueDate: &dueDate}
	if err := tasks.Create(t.Context(), task); err != nil {
		t.Fatalf("seeding task: %v", err)
	}

	w := doAuthedJSONRequest(router, http.MethodPatch, "/workspaces/"+workspaceID.String()+"/tasks/"+task.ID.String(), bearerToken(t, manager, callerID), map[string]any{
		"clearSprint":  true,
		"clearDueDate": true,
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
	var body struct {
		Data taskResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Data.SprintID != nil {
		t.Errorf("sprintId = %v, want nil after clearSprint", body.Data.SprintID)
	}
	if body.Data.DueDate != nil {
		t.Errorf("dueDate = %v, want nil after clearDueDate", body.Data.DueDate)
	}
}

func TestTaskHandler_Delete_Success(t *testing.T) {
	h, tasks, projects, members, manager := newTestTaskHandler()
	router := newTaskTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)
	project := projects.seedProject(workspaceID, "Launch")
	task := &domain.Task{WorkspaceID: workspaceID, ProjectID: project.ID, Title: "Ship it"}
	if err := tasks.Create(t.Context(), task); err != nil {
		t.Fatalf("seeding task: %v", err)
	}

	deleteW := doAuthedJSONRequest(router, http.MethodDelete, "/workspaces/"+workspaceID.String()+"/tasks/"+task.ID.String(), bearerToken(t, manager, callerID), nil)
	if deleteW.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", deleteW.Code, http.StatusOK, deleteW.Body.String())
	}

	getW := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceID.String()+"/tasks/"+task.ID.String(), bearerToken(t, manager, callerID), nil)
	if getW.Code != http.StatusNotFound {
		t.Errorf("status after delete = %d, want %d", getW.Code, http.StatusNotFound)
	}
}

func TestTaskHandler_List_FiltersByStatus(t *testing.T) {
	h, _, projects, members, manager := newTestTaskHandler()
	router := newTaskTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)
	project := projects.seedProject(workspaceID, "Launch")

	doAuthedJSONRequest(router, http.MethodPost, "/workspaces/"+workspaceID.String()+"/tasks", bearerToken(t, manager, callerID), map[string]any{
		"projectId": project.ID.String(),
		"title":     "Ship it",
	})

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceID.String()+"/tasks?status=done", bearerToken(t, manager, callerID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
	var body struct {
		Data []taskResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(body.Data) != 0 {
		t.Errorf("filtering by status=done: got %d tasks, want 0 (the seeded task is still todo)", len(body.Data))
	}
}
