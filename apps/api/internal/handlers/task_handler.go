package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
)

// TaskHandler serves the /workspaces/{workspaceId}/tasks routes (see
// docs/api.md). Every route here must have middleware.Auth applied — they're
// all membership-gated, so there's no anonymous access path. It reuses
// requireUserID/parseUUIDParam from workspace_handler.go rather than
// redefining them.
type TaskHandler struct {
	taskUsecase *usecases.TaskUsecase
}

// NewTaskHandler constructs a TaskHandler.
func NewTaskHandler(taskUsecase *usecases.TaskUsecase) *TaskHandler {
	return &TaskHandler{taskUsecase: taskUsecase}
}

type createTaskRequest struct {
	ProjectID  string  `json:"projectId" binding:"required"`
	Title      string  `json:"title" binding:"required"`
	AssigneeID *string `json:"assigneeId"`
}

// updateTaskRequest's pointer fields use nil-means-omitted for a partial
// PATCH (same convention as updateWorkspaceRequest/updateProjectRequest).
// ClearAssignee is a separate explicit flag because a raw JSON null for
// assigneeId would bind to the same Go nil as the field being omitted
// entirely — see usecases.UpdateTaskInput.ClearAssignee.
type updateTaskRequest struct {
	ProjectID     *string `json:"projectId"`
	Title         *string `json:"title"`
	Status        *string `json:"status"`
	AssigneeID    *string `json:"assigneeId"`
	ClearAssignee bool    `json:"clearAssignee"`
}

type taskResponse struct {
	ID          string  `json:"id"`
	WorkspaceID string  `json:"workspaceId"`
	ProjectID   string  `json:"projectId"`
	Title       string  `json:"title"`
	Status      string  `json:"status"`
	AssigneeID  *string `json:"assigneeId"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func toTaskResponse(t *domain.Task) taskResponse {
	resp := taskResponse{
		ID:          t.ID.String(),
		WorkspaceID: t.WorkspaceID.String(),
		ProjectID:   t.ProjectID.String(),
		Title:       t.Title,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
	if t.AssigneeID != nil {
		id := t.AssigneeID.String()
		resp.AssigneeID = &id
	}
	return resp
}

// Create handles POST /workspaces/{workspaceId}/tasks.
func (h *TaskHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}

	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		handleError(c, apperrors.Validation("projectId", "invalid id"))
		return
	}

	input := usecases.CreateTaskInput{ProjectID: projectID, Title: req.Title}
	if req.AssigneeID != nil {
		assigneeID, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			handleError(c, apperrors.Validation("assigneeId", "invalid id"))
			return
		}
		input.AssigneeID = &assigneeID
	}

	task, err := h.taskUsecase.Create(c.Request.Context(), userID, workspaceID, input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": toTaskResponse(task)})
}

// List handles GET /workspaces/{workspaceId}/tasks, optionally filtered by
// the ?projectId=/?status=/?assigneeId= query params.
func (h *TaskHandler) List(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}

	var filter domain.TaskFilter
	if raw := c.Query("projectId"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			handleError(c, apperrors.Validation("projectId", "invalid id"))
			return
		}
		filter.ProjectID = &id
	}
	if raw := c.Query("status"); raw != "" {
		status := domain.TaskStatus(raw)
		filter.Status = &status
	}
	if raw := c.Query("assigneeId"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			handleError(c, apperrors.Validation("assigneeId", "invalid id"))
			return
		}
		filter.AssigneeID = &id
	}

	tasks, err := h.taskUsecase.ListByWorkspaceID(c.Request.Context(), userID, workspaceID, filter)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = toTaskResponse(t)
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Get handles GET /workspaces/{workspaceId}/tasks/{taskId}.
func (h *TaskHandler) Get(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}
	taskID, ok := parseUUIDParam(c, "taskId")
	if !ok {
		return
	}

	task, err := h.taskUsecase.GetByID(c.Request.Context(), userID, workspaceID, taskID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toTaskResponse(task)})
}

// Update handles PATCH /workspaces/{workspaceId}/tasks/{taskId}.
func (h *TaskHandler) Update(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}
	taskID, ok := parseUUIDParam(c, "taskId")
	if !ok {
		return
	}

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	input := usecases.UpdateTaskInput{Title: req.Title, ClearAssignee: req.ClearAssignee}
	if req.ProjectID != nil {
		id, err := uuid.Parse(*req.ProjectID)
		if err != nil {
			handleError(c, apperrors.Validation("projectId", "invalid id"))
			return
		}
		input.ProjectID = &id
	}
	if req.Status != nil {
		status := domain.TaskStatus(*req.Status)
		input.Status = &status
	}
	if req.AssigneeID != nil {
		id, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			handleError(c, apperrors.Validation("assigneeId", "invalid id"))
			return
		}
		input.AssigneeID = &id
	}

	task, err := h.taskUsecase.Update(c.Request.Context(), userID, workspaceID, taskID, input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toTaskResponse(task)})
}

// Delete handles DELETE /workspaces/{workspaceId}/tasks/{taskId}.
func (h *TaskHandler) Delete(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}
	taskID, ok := parseUUIDParam(c, "taskId")
	if !ok {
		return
	}

	if err := h.taskUsecase.Delete(c.Request.Context(), userID, workspaceID, taskID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
