package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
)

// ProjectHandler serves the /workspaces/{workspaceId}/projects routes (see
// docs/api.md). Every route here must have middleware.Auth applied — they're
// all membership-gated, so there's no anonymous access path.
type ProjectHandler struct {
	projectUsecase *usecases.ProjectUsecase
}

// NewProjectHandler constructs a ProjectHandler.
func NewProjectHandler(projectUsecase *usecases.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{projectUsecase: projectUsecase}
}

type createProjectRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// updateProjectRequest's fields are pointers so ShouldBindJSON can tell
// "field omitted" (nil, leave unchanged) from "field explicitly cleared"
// for a partial PATCH — same convention as updateWorkspaceRequest.
type updateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type projectResponse struct {
	ID          string  `json:"id"`
	WorkspaceID string  `json:"workspaceId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func toProjectResponse(p *domain.Project) projectResponse {
	return projectResponse{
		ID:          p.ID.String(),
		WorkspaceID: p.WorkspaceID.String(),
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}

// Create handles POST /workspaces/{workspaceId}/projects.
func (h *ProjectHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}

	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	project, err := h.projectUsecase.Create(c.Request.Context(), userID, workspaceID, req.Name, req.Description)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": toProjectResponse(project)})
}

// List handles GET /workspaces/{workspaceId}/projects.
func (h *ProjectHandler) List(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}

	projects, err := h.projectUsecase.ListByWorkspaceID(c.Request.Context(), userID, workspaceID)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := make([]projectResponse, len(projects))
	for i, p := range projects {
		resp[i] = toProjectResponse(p)
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Get handles GET /workspaces/{workspaceId}/projects/{id}.
func (h *ProjectHandler) Get(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}
	projectID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	project, err := h.projectUsecase.GetByID(c.Request.Context(), userID, workspaceID, projectID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toProjectResponse(project)})
}

// Update handles PATCH /workspaces/{workspaceId}/projects/{id}.
func (h *ProjectHandler) Update(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}
	projectID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req updateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	project, err := h.projectUsecase.Update(c.Request.Context(), userID, workspaceID, projectID, usecases.UpdateProjectInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toProjectResponse(project)})
}

// Delete handles DELETE /workspaces/{workspaceId}/projects/{id}.
func (h *ProjectHandler) Delete(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "workspaceId")
	if !ok {
		return
	}
	projectID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.projectUsecase.Delete(c.Request.Context(), userID, workspaceID, projectID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}
