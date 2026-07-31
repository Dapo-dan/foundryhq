package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
)

func newTestProjectHandler() (*ProjectHandler, *fakeProjectRepo, *fakeWorkspaceMemberRepo, *jwt.Manager) {
	projects := newFakeProjectRepo()
	members := newFakeWorkspaceMemberRepo()
	usecase := usecases.NewProjectUsecase(projects, members)
	return NewProjectHandler(usecase), projects, members, newTestJWTManager()
}

func newProjectTestRouter(manager *jwt.Manager, h *ProjectHandler) *gin.Engine {
	return newProtectedTestRouter(manager, func(protected *gin.RouterGroup) {
		group := protected.Group("/workspaces/:workspaceId/projects")
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
	})
}

func TestProjectHandler_Create_Success(t *testing.T) {
	h, _, members, manager := newTestProjectHandler()
	router := newProjectTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleOwner)

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces/"+workspaceID.String()+"/projects", bearerToken(t, manager, callerID), map[string]string{
		"name": "Launch",
	})

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusCreated, w.Body.String())
	}
	var body struct {
		Data projectResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Data.WorkspaceID != workspaceID.String() {
		t.Errorf("workspaceId = %q, want %q", body.Data.WorkspaceID, workspaceID.String())
	}
}

func TestProjectHandler_Create_ForbiddenForNonMember(t *testing.T) {
	h, _, _, manager := newTestProjectHandler()
	router := newProjectTestRouter(manager, h)
	workspaceID := uuid.New()

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces/"+workspaceID.String()+"/projects", bearerToken(t, manager, uuid.New()), map[string]string{
		"name": "Launch",
	})

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusForbidden, w.Body.String())
	}
	assertErrorCode(t, w, apperrors.CodeForbidden)
}

func TestProjectHandler_Get_CrossWorkspaceIsNotFound(t *testing.T) {
	h, projects, members, manager := newTestProjectHandler()
	router := newProjectTestRouter(manager, h)
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	members.seedMember(workspaceA, callerID, domain.RoleMember)
	projectInB := projects.seedProject(workspaceB, "Other workspace's project")

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceA.String()+"/projects/"+projectInB.ID.String(), bearerToken(t, manager, callerID), nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body = %s — a project must not be reachable through a different workspace's path", w.Code, http.StatusNotFound, w.Body.String())
	}
}

func TestProjectHandler_List_Success(t *testing.T) {
	h, projects, members, manager := newTestProjectHandler()
	router := newProjectTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)
	projects.seedProject(workspaceID, "Launch")

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceID.String()+"/projects", bearerToken(t, manager, callerID), nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
	var body struct {
		Data []projectResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(body.Data) != 1 {
		t.Fatalf("expected 1 project, got %d", len(body.Data))
	}
}

// TestProjectHandler_Delete_Success only checks the HTTP contract (200, then
// a 404 on a subsequent Get) — the project→tasks cascade itself lives inside
// the real Postgres repository's transaction (domain.ProjectRepository has
// no notion of tasks), and is already verified against live Postgres by
// TestProjectRepository_Delete_CascadesToTasks; disconnected fakes can't
// meaningfully re-exercise that at this layer.
func TestProjectHandler_Delete_Success(t *testing.T) {
	h, projects, members, manager := newTestProjectHandler()
	router := newProjectTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleOwner)
	project := projects.seedProject(workspaceID, "Launch")

	deleteW := doAuthedJSONRequest(router, http.MethodDelete, "/workspaces/"+workspaceID.String()+"/projects/"+project.ID.String(), bearerToken(t, manager, callerID), nil)
	if deleteW.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", deleteW.Code, http.StatusOK, deleteW.Body.String())
	}

	getW := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceID.String()+"/projects/"+project.ID.String(), bearerToken(t, manager, callerID), nil)
	if getW.Code != http.StatusNotFound {
		t.Errorf("status after delete = %d, want %d — a soft-deleted project should read as gone", getW.Code, http.StatusNotFound)
	}
}
