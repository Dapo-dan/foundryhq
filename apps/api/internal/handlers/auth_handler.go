package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
)

const refreshTokenCookieName = "refresh_token"

// refreshTokenCookiePath scopes the cookie to only the two endpoints that
// need it, rather than sending it on every request to the API.
const refreshTokenCookiePath = "/auth"

// AuthHandler serves POST /auth/register, /auth/login, /auth/refresh, and
// /auth/logout (see docs/api.md). These routes must not have
// middleware.Auth applied — they're how a caller obtains a session in the
// first place.
type AuthHandler struct {
	authUsecase   *usecases.AuthUsecase
	secureCookies bool
}

// NewAuthHandler constructs an AuthHandler. secureCookies should be true in
// production (HTTPS) and false for local HTTP dev — see
// adr/0004-jwt-access-refresh-tokens.md.
func NewAuthHandler(authUsecase *usecases.AuthUsecase, secureCookies bool) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, secureCookies: secureCookies}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// userResponse/authResponse use camelCase json tags — not the snake_case
// convention in .ai/documentation/api.md — because they must match
// packages/shared-types/src/auth.ts's AuthSession shape, which the already-
// built sign-in/sign-up screens consume directly. Matching the real
// frontend contract takes priority over the stale doc convention here.
type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type authResponse struct {
	User        userResponse `json:"user"`
	AccessToken string       `json:"accessToken"`
}

func toAuthResponse(result *usecases.AuthResult) authResponse {
	return authResponse{
		User: userResponse{
			ID:    result.User.ID.String(),
			Email: result.User.Email,
		},
		AccessToken: result.AccessToken,
	}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	result, err := h.authUsecase.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	h.setRefreshCookie(c, result.RefreshToken, result.RefreshExpiresAt)
	c.JSON(http.StatusCreated, gin.H{"data": toAuthResponse(result)})
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.Validation("", err.Error()))
		return
	}

	result, err := h.authUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	h.setRefreshCookie(c, result.RefreshToken, result.RefreshExpiresAt)
	c.JSON(http.StatusOK, gin.H{"data": toAuthResponse(result)})
}

// Refresh handles POST /auth/refresh. The refresh token is read from the
// httpOnly cookie, never the request body — the browser attaches it
// automatically (see api-client.ts's withCredentials: true), and JS can't
// read it either way.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || refreshToken == "" {
		handleError(c, apperrors.Unauthorized("missing refresh token"))
		return
	}

	result, err := h.authUsecase.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		handleError(c, err)
		return
	}

	h.setRefreshCookie(c, result.RefreshToken, result.RefreshExpiresAt)
	c.JSON(http.StatusOK, gin.H{"data": toAuthResponse(result)})
}

// Logout handles POST /auth/logout. Missing/already-invalid cookies aren't
// an error — the caller ends up logged out either way.
func (h *AuthHandler) Logout(c *gin.Context) {
	if refreshToken, err := c.Cookie(refreshTokenCookieName); err == nil && refreshToken != "" {
		if err := h.authUsecase.Logout(c.Request.Context(), refreshToken); err != nil {
			handleError(c, err)
			return
		}
	}

	h.clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshTokenCookieName, token, maxAge, refreshTokenCookiePath, "", h.secureCookies, true)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshTokenCookieName, "", -1, refreshTokenCookiePath, "", h.secureCookies, true)
}
