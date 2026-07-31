package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
)

func newTestSprintHandler() (*SprintHandler, *fakeSprintRepo, *fakeTaskRepo, *fakeWorkspaceMemberRepo, *jwt.Manager) {
	sprints := newFakeSprintRepo()
	tasks := newFakeTaskRepo()
	members := newFakeWorkspaceMemberRepo()
	usecase := usecases.NewSprintUsecase(sprints, tasks, members)
	return NewSprintHandler(usecase), sprints, tasks, members, newTestJWTManager()
}

func newSprintTestRouter(manager *jwt.Manager, h *SprintHandler) *gin.Engine {
	return newProtectedTestRouter(manager, func(protected *gin.RouterGroup) {
		group := protected.Group("/workspaces/:workspaceId/sprints")
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:sprintId", h.Get)
		group.GET("/:sprintId/velocity", h.Velocity)
	})
}

func TestSprintHandler_Create_Success(t *testing.T) {
	h, _, _, members, manager := newTestSprintHandler()
	router := newSprintTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces/"+workspaceID.String()+"/sprints", bearerToken(t, manager, callerID), map[string]string{
		"name":      "Sprint 1",
		"startDate": "2026-08-01",
		"endDate":   "2026-08-14",
	})

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusCreated, w.Body.String())
	}
}

func TestSprintHandler_Create_EndBeforeStart(t *testing.T) {
	h, _, _, members, manager := newTestSprintHandler()
	router := newSprintTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces/"+workspaceID.String()+"/sprints", bearerToken(t, manager, callerID), map[string]string{
		"name":      "Backwards sprint",
		"startDate": "2026-08-14",
		"endDate":   "2026-08-01",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

func TestSprintHandler_Get_ReturnsSprintWithTasks(t *testing.T) {
	h, sprints, tasks, members, manager := newTestSprintHandler()
	router := newSprintTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprint := sprints.seedSprint(workspaceID, "Sprint 1", start, end)
	task := &domain.Task{WorkspaceID: workspaceID, ProjectID: uuid.New(), Title: "Ship it", SprintID: &sprint.ID}
	if err := tasks.Create(t.Context(), task); err != nil {
		t.Fatalf("seeding task: %v", err)
	}

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceID.String()+"/sprints/"+sprint.ID.String(), bearerToken(t, manager, callerID), nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
	var body struct {
		Data sprintWithTasksResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(body.Data.Tasks) != 1 {
		t.Fatalf("expected 1 task in sprint detail, got %d", len(body.Data.Tasks))
	}
}

func TestSprintHandler_Get_CrossWorkspaceIsNotFound(t *testing.T) {
	h, sprints, _, members, manager := newTestSprintHandler()
	router := newSprintTestRouter(manager, h)
	workspaceA, workspaceB, callerID := uuid.New(), uuid.New(), uuid.New()
	members.seedMember(workspaceA, callerID, domain.RoleMember)
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprintInB := sprints.seedSprint(workspaceB, "Other workspace's sprint", start, end)

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceA.String()+"/sprints/"+sprintInB.ID.String(), bearerToken(t, manager, callerID), nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body = %s — a sprint must not be reachable through a different workspace's path", w.Code, http.StatusNotFound, w.Body.String())
	}
}

// TestSprintHandler_Velocity_SumsOnlyDoneWithinRange exercises the exact
// business rule from .ai/business-analysis/acceptance-criteria.md's AC for
// REQ-05: velocity is the sum of story points for tasks done within the
// sprint's date range — done-before, done-after, and not-done tasks must
// all be excluded, and a NULL story_points value must not break the sum.
func TestSprintHandler_Velocity_SumsOnlyDoneWithinRange(t *testing.T) {
	h, sprints, tasks, members, manager := newTestSprintHandler()
	router := newSprintTestRouter(manager, h)
	workspaceID, callerID := uuid.New(), uuid.New()
	members.seedMember(workspaceID, callerID, domain.RoleMember)

	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprint := sprints.seedSprint(workspaceID, "Sprint 1", start, end)

	inRangePoints := 5
	doneInRange := &domain.Task{WorkspaceID: workspaceID, ProjectID: uuid.New(), Title: "Done in range", SprintID: &sprint.ID, Status: domain.StatusDone, StoryPoints: &inRangePoints}
	if err := tasks.Create(t.Context(), doneInRange); err != nil {
		t.Fatalf("seeding task: %v", err)
	}
	doneInRange.UpdatedAt = time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)

	notDonePoints := 8
	notDone := &domain.Task{WorkspaceID: workspaceID, ProjectID: uuid.New(), Title: "Still in progress", SprintID: &sprint.ID, Status: domain.StatusInProgress, StoryPoints: &notDonePoints}
	if err := tasks.Create(t.Context(), notDone); err != nil {
		t.Fatalf("seeding task: %v", err)
	}

	doneAfterClosePoints := 13
	doneAfterClose := &domain.Task{WorkspaceID: workspaceID, ProjectID: uuid.New(), Title: "Done after sprint closed", SprintID: &sprint.ID, Status: domain.StatusDone, StoryPoints: &doneAfterClosePoints}
	if err := tasks.Create(t.Context(), doneAfterClose); err != nil {
		t.Fatalf("seeding task: %v", err)
	}
	doneAfterClose.UpdatedAt = time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)

	doneNullPoints := &domain.Task{WorkspaceID: workspaceID, ProjectID: uuid.New(), Title: "Done, no points set", SprintID: &sprint.ID, Status: domain.StatusDone}
	if err := tasks.Create(t.Context(), doneNullPoints); err != nil {
		t.Fatalf("seeding task: %v", err)
	}
	doneNullPoints.UpdatedAt = time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC)

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceID.String()+"/sprints/"+sprint.ID.String()+"/velocity", bearerToken(t, manager, callerID), nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
	var body struct {
		Data struct {
			Velocity int `json:"velocity"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	// Only doneInRange's 5 points should count — notDone isn't done,
	// doneAfterClose finished after the sprint's end_date, and
	// doneNullPoints has no story_points to add (COALESCE(...,0), not NULL).
	if body.Data.Velocity != 5 {
		t.Errorf("velocity = %d, want %d", body.Data.Velocity, 5)
	}
}

func TestSprintHandler_Velocity_ForbiddenForNonMember(t *testing.T) {
	h, sprints, _, _, manager := newTestSprintHandler()
	router := newSprintTestRouter(manager, h)
	workspaceID := uuid.New()
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	sprint := sprints.seedSprint(workspaceID, "Sprint 1", start, end)

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspaceID.String()+"/sprints/"+sprint.ID.String()+"/velocity", bearerToken(t, manager, uuid.New()), nil)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusForbidden, w.Body.String())
	}
}
