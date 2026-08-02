package handlers

import (
	"context"
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

func newTestWorkspaceHandler() (*WorkspaceHandler, *fakeWorkspaceRepo, *fakeWorkspaceMemberRepo, *jwt.Manager) {
	workspaces := newFakeWorkspaceRepo()
	members := newFakeWorkspaceMemberRepo()
	usecase := usecases.NewWorkspaceUsecase(
		workspaces, members, newFakeUserRepo(),
		newFakeInviteTokenRepo(), &fakeMailer{}, "http://localhost:5173",
	)
	return NewWorkspaceHandler(usecase), workspaces, members, newTestJWTManager()
}

func newWorkspaceTestRouter(manager *jwt.Manager, h *WorkspaceHandler) *gin.Engine {
	return newProtectedTestRouter(manager, func(g *gin.RouterGroup) {
		g.POST("/workspaces", h.Create)
		g.GET("/workspaces", h.List)
		g.GET("/workspaces/:workspaceId", h.Get)
		g.PATCH("/workspaces/:workspaceId", h.Update)
		g.GET("/workspaces/:workspaceId/members", h.ListMembers)
		g.POST("/workspaces/:workspaceId/members/invite", h.Invite)
		g.PATCH("/workspaces/:workspaceId/members/:memberId", h.UpdateMemberRole)
	})
}

func TestWorkspaceHandler_Create_Success(t *testing.T) {
	h, _, _, manager := newTestWorkspaceHandler()
	router := newWorkspaceTestRouter(manager, h)
	userID := uuid.New()

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces", bearerToken(t, manager, userID), map[string]string{
		"name": "Acme Inc.",
	})

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusCreated, w.Body.String())
	}
	var body struct {
		Data workspaceResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Data.Slug != "acme-inc" {
		t.Errorf("slug = %q, want %q", body.Data.Slug, "acme-inc")
	}
}

func TestWorkspaceHandler_Create_NoAuthHeader(t *testing.T) {
	h, _, _, manager := newTestWorkspaceHandler()
	router := newWorkspaceTestRouter(manager, h)

	req, w := jsonRequestNoAuth(http.MethodPost, "/workspaces", map[string]string{"name": "Acme Inc."})
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d — middleware.Auth should reject a missing token", w.Code, http.StatusUnauthorized)
	}
}

func TestWorkspaceHandler_Get_ForbiddenForNonMember(t *testing.T) {
	h, workspaces, _, manager := newTestWorkspaceHandler()
	router := newWorkspaceTestRouter(manager, h)
	ownerID, strangerID := uuid.New(), uuid.New()

	workspace := seedWorkspaceWithOwner(workspaces, "Acme Inc.", ownerID)

	w := doAuthedJSONRequest(router, http.MethodGet, "/workspaces/"+workspace.ID.String(), bearerToken(t, manager, strangerID), nil)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusForbidden, w.Body.String())
	}
	assertErrorCode(t, w, apperrors.CodeForbidden)
}

func TestWorkspaceHandler_Update_ForbiddenForNonOwner(t *testing.T) {
	h, workspaces, members, manager := newTestWorkspaceHandler()
	router := newWorkspaceTestRouter(manager, h)
	ownerID, memberID := uuid.New(), uuid.New()

	workspace := seedWorkspaceWithOwner(workspaces, "Acme Inc.", ownerID)
	members.seedMember(workspace.ID, ownerID, domain.RoleOwner)
	members.seedMember(workspace.ID, memberID, domain.RoleMember)

	w := doAuthedJSONRequest(router, http.MethodPatch, "/workspaces/"+workspace.ID.String(), bearerToken(t, manager, memberID), map[string]string{
		"name": "Hijacked",
	})

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d, body = %s — a plain member must not be able to rename the workspace", w.Code, http.StatusForbidden, w.Body.String())
	}
}

func TestWorkspaceHandler_Update_OwnerSucceeds(t *testing.T) {
	h, workspaces, members, manager := newTestWorkspaceHandler()
	router := newWorkspaceTestRouter(manager, h)
	ownerID := uuid.New()

	workspace := seedWorkspaceWithOwner(workspaces, "Acme Inc.", ownerID)
	members.seedMember(workspace.ID, ownerID, domain.RoleOwner)

	w := doAuthedJSONRequest(router, http.MethodPatch, "/workspaces/"+workspace.ID.String(), bearerToken(t, manager, ownerID), map[string]string{
		"name": "Acme Corp",
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestWorkspaceHandler_UpdateMemberRole_ForbiddenForNonOwner(t *testing.T) {
	h, workspaces, members, manager := newTestWorkspaceHandler()
	router := newWorkspaceTestRouter(manager, h)
	ownerID, memberID := uuid.New(), uuid.New()

	workspace := seedWorkspaceWithOwner(workspaces, "Acme Inc.", ownerID)
	ownerMember := members.seedMember(workspace.ID, ownerID, domain.RoleOwner)
	members.seedMember(workspace.ID, memberID, domain.RoleMember)

	w := doAuthedJSONRequest(router, http.MethodPatch,
		"/workspaces/"+workspace.ID.String()+"/members/"+ownerMember.ID.String(),
		bearerToken(t, manager, memberID),
		map[string]string{"role": "member"},
	)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d, body = %s — a plain member must not be able to change roles", w.Code, http.StatusForbidden, w.Body.String())
	}
}

func TestWorkspaceHandler_Invite_Success(t *testing.T) {
	h, workspaces, members, manager := newTestWorkspaceHandler()
	router := newWorkspaceTestRouter(manager, h)
	ownerID := uuid.New()

	workspace := seedWorkspaceWithOwner(workspaces, "Acme Inc.", ownerID)
	members.seedMember(workspace.ID, ownerID, domain.RoleOwner)

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces/"+workspace.ID.String()+"/members/invite", bearerToken(t, manager, ownerID), map[string]string{
		"email": "new@example.com",
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestWorkspaceHandler_Create_ValidationError(t *testing.T) {
	h, _, _, manager := newTestWorkspaceHandler()
	router := newWorkspaceTestRouter(manager, h)

	w := doAuthedJSONRequest(router, http.MethodPost, "/workspaces", bearerToken(t, manager, uuid.New()), map[string]string{
		"name": "   ",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
	assertErrorCode(t, w, apperrors.CodeValidation)
}

// seedWorkspaceWithOwner inserts a workspace directly (bypassing the
// usecase, which handler tests don't need) and hands back the created
// domain.Workspace so tests have a real ID to address.
func seedWorkspaceWithOwner(repo *fakeWorkspaceRepo, name string, ownerID uuid.UUID) *domain.Workspace {
	workspace := &domain.Workspace{Name: name, Slug: name}
	owner := &domain.WorkspaceMember{UserID: ownerID, Role: domain.RoleOwner}
	_ = repo.Create(context.Background(), workspace, owner)
	return workspace
}
