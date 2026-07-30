package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/middleware"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
)

// WorkspaceHandler serves the /workspaces routes (see docs/api.md). Every
// route here must have middleware.Auth applied — they're all
// membership-gated, so there's no anonymous access path.
type WorkspaceHandler struct {
	workspaceUsecase *usecases.WorkspaceUsecase
}

// NewWorkspaceHandler constructs a WorkspaceHandler.
func NewWorkspaceHandler(workspaceUsecase *usecases.WorkspaceUsecase) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceUsecase: workspaceUsecase}
}

type createWorkspaceRequest struct {
	Name string `json:"name" binding:"required"`
}

// updateWorkspaceRequest's fields are pointers so ShouldBindJSON can tell
// "field omitted" (nil, leave unchanged) from "field explicitly cleared"
// (empty string) for a partial PATCH.
type updateWorkspaceRequest struct {
	Name    *string `json:"name"`
	Slug    *string `json:"slug"`
	LogoURL *string `json:"logoUrl"`
}

type inviteMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type updateMemberRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type workspaceResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	LogoURL   string `json:"logoUrl"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func toWorkspaceResponse(w *domain.Workspace) workspaceResponse {
	return workspaceResponse{
		ID:        w.ID.String(),
		Name:      w.Name,
		Slug:      w.Slug,
		LogoURL:   w.LogoURL,
		CreatedAt: w.CreatedAt.Format(time.RFC3339),
		UpdatedAt: w.UpdatedAt.Format(time.RFC3339),
	}
}

type workspaceMemberResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"userId"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	InvitedAt string  `json:"invitedAt"`
	JoinedAt  *string `json:"joinedAt"`
}

func toWorkspaceMemberResponse(m *domain.WorkspaceMemberWithUser) workspaceMemberResponse {
	resp := workspaceMemberResponse{
		ID:        m.ID.String(),
		UserID:    m.UserID.String(),
		Email:     m.Email,
		Role:      string(m.Role),
		InvitedAt: m.InvitedAt.Format(time.RFC3339),
	}
	if m.JoinedAt != nil {
		joinedAt := m.JoinedAt.Format(time.RFC3339)
		resp.JoinedAt = &joinedAt
	}
	return resp
}

// Create handles POST /workspaces.
func (h *WorkspaceHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	var req createWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	workspace, err := h.workspaceUsecase.Create(c.Request.Context(), userID, req.Name)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": toWorkspaceResponse(workspace)})
}

// List handles GET /workspaces — every workspace the caller belongs to.
func (h *WorkspaceHandler) List(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	workspaces, err := h.workspaceUsecase.ListForUser(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := make([]workspaceResponse, len(workspaces))
	for i, w := range workspaces {
		resp[i] = toWorkspaceResponse(w)
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Get handles GET /workspaces/{id}.
func (h *WorkspaceHandler) Get(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	workspace, err := h.workspaceUsecase.GetByID(c.Request.Context(), userID, workspaceID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toWorkspaceResponse(workspace)})
}

// Update handles PATCH /workspaces/{id}.
func (h *WorkspaceHandler) Update(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req updateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	workspace, err := h.workspaceUsecase.Update(c.Request.Context(), userID, workspaceID, usecases.UpdateWorkspaceInput{
		Name:    req.Name,
		Slug:    req.Slug,
		LogoURL: req.LogoURL,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toWorkspaceResponse(workspace)})
}

// ListMembers handles GET /workspaces/{id}/members.
func (h *WorkspaceHandler) ListMembers(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	members, err := h.workspaceUsecase.ListMembers(c.Request.Context(), userID, workspaceID)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := make([]workspaceMemberResponse, len(members))
	for i, m := range members {
		resp[i] = toWorkspaceMemberResponse(m)
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Invite handles POST /workspaces/{id}/members/invite.
func (h *WorkspaceHandler) Invite(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req inviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	if _, err := h.workspaceUsecase.Invite(c.Request.Context(), userID, workspaceID, req.Email); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}

// UpdateMemberRole handles PATCH /workspaces/{id}/members/{memberId}.
func (h *WorkspaceHandler) UpdateMemberRole(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	memberID, ok := parseUUIDParam(c, "memberId")
	if !ok {
		return
	}

	var req updateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	err := h.workspaceUsecase.UpdateMemberRole(c.Request.Context(), userID, workspaceID, memberID, domain.WorkspaceRole(req.Role))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}

// requireUserID reads the caller's user ID set by middleware.Auth. false
// means Auth didn't run for this request — a route registered without the
// middleware, which is a programmer error, not a real 401.
func requireUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		handleError(c, apperrors.Internal(errors.New("missing user id in context: route registered without middleware.Auth")))
		return uuid.UUID{}, false
	}
	return userID, true
}

// parseUUIDParam parses the named path param as a UUID, responding with a
// validation error and returning false if it isn't one.
func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		handleError(c, apperrors.Validation(name, "invalid id"))
		return uuid.UUID{}, false
	}
	return id, true
}
