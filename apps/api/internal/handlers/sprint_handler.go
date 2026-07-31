package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
)

// dateOnlyLayout parses/formats sprints.start_date/end_date — plain
// calendar dates, no time-of-day component.
const dateOnlyLayout = "2006-01-02"

// SprintHandler serves the /workspaces/{workspaceId}/sprints routes (see
// docs/api.md). Every route here must have middleware.Auth applied —
// they're all membership-gated, so there's no anonymous access path. It
// reuses requireUserID/parseUUIDParam from workspace_handler.go rather than
// redefining them.
type SprintHandler struct {
	sprintUsecase *usecases.SprintUsecase
}

// NewSprintHandler constructs a SprintHandler.
func NewSprintHandler(sprintUsecase *usecases.SprintUsecase) *SprintHandler {
	return &SprintHandler{sprintUsecase: sprintUsecase}
}

type createSprintRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
}

type sprintResponse struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func toSprintResponse(s *domain.Sprint) sprintResponse {
	return sprintResponse{
		ID:          s.ID.String(),
		WorkspaceID: s.WorkspaceID.String(),
		Name:        s.Name,
		StartDate:   s.StartDate.Format(dateOnlyLayout),
		EndDate:     s.EndDate.Format(dateOnlyLayout),
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}
}

type sprintWithTasksResponse struct {
	Sprint sprintResponse `json:"sprint"`
	Tasks  []taskResponse `json:"tasks"`
}

// Create handles POST /workspaces/{workspaceId}/sprints.
func (h *SprintHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}

	var req createSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	startDate, err := time.Parse(dateOnlyLayout, req.StartDate)
	if err != nil {
		handleError(c, apperrors.Validation("startDate", "must be a date in YYYY-MM-DD format"))
		return
	}
	endDate, err := time.Parse(dateOnlyLayout, req.EndDate)
	if err != nil {
		handleError(c, apperrors.Validation("endDate", "must be a date in YYYY-MM-DD format"))
		return
	}

	sprint, err := h.sprintUsecase.Create(c.Request.Context(), userID, workspaceID, req.Name, startDate, endDate)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": toSprintResponse(sprint)})
}

// List handles GET /workspaces/{workspaceId}/sprints.
func (h *SprintHandler) List(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}

	sprints, err := h.sprintUsecase.ListForWorkspace(c.Request.Context(), userID, workspaceID)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := make([]sprintResponse, len(sprints))
	for i, s := range sprints {
		resp[i] = toSprintResponse(s)
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Get handles GET /workspaces/{workspaceId}/sprints/{sprintId} — the sprint
// plus its tasks. Grouping by status is left to the caller (the web client
// reuses the same Kanban board component the Tasks page uses).
func (h *SprintHandler) Get(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}
	sprintID, ok := parseUUIDParam(c, "sprintId")
	if !ok {
		return
	}

	result, err := h.sprintUsecase.GetByID(c.Request.Context(), userID, workspaceID, sprintID)
	if err != nil {
		handleError(c, err)
		return
	}

	tasks := make([]taskResponse, len(result.Tasks))
	for i, t := range result.Tasks {
		tasks[i] = toTaskResponse(t)
	}
	c.JSON(http.StatusOK, gin.H{"data": sprintWithTasksResponse{
		Sprint: toSprintResponse(result.Sprint),
		Tasks:  tasks,
	}})
}

// Velocity handles GET /workspaces/{workspaceId}/sprints/{sprintId}/velocity.
func (h *SprintHandler) Velocity(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}
	sprintID, ok := parseUUIDParam(c, "sprintId")
	if !ok {
		return
	}

	velocity, err := h.sprintUsecase.GetVelocity(c.Request.Context(), userID, workspaceID, sprintID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"velocity": velocity}})
}
