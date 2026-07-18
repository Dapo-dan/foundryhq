package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
)

// fakeUserRepo/fakeRefreshTokenRepo are minimal in-memory doubles, kept
// local to this package (mirroring, not sharing, usecases' own test
// doubles) so handler tests stay decoupled from usecase internals.

type fakeUserRepo struct {
	byEmail map[string]*domain.User
	byID    map[uuid.UUID]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byEmail: map[string]*domain.User{}, byID: map[uuid.UUID]*domain.User{}}
}

func (r *fakeUserRepo) Create(_ context.Context, user *domain.User) error {
	if _, exists := r.byEmail[user.Email]; exists {
		return domain.ErrEmailAlreadyExists
	}
	user.ID = uuid.New()
	r.byEmail[user.Email] = user
	r.byID[user.ID] = user
	return nil
}

func (r *fakeUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (r *fakeUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	u, ok := r.byEmail[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

type fakeRefreshTokenRepo struct {
	byHash map[string]*domain.RefreshToken
}

func newFakeRefreshTokenRepo() *fakeRefreshTokenRepo {
	return &fakeRefreshTokenRepo{byHash: map[string]*domain.RefreshToken{}}
}

func (r *fakeRefreshTokenRepo) Create(_ context.Context, token *domain.RefreshToken) error {
	token.ID = uuid.New()
	r.byHash[token.TokenHash] = token
	return nil
}

func (r *fakeRefreshTokenRepo) GetByTokenHash(_ context.Context, tokenHash string) (*domain.RefreshToken, error) {
	t, ok := r.byHash[tokenHash]
	if !ok {
		return nil, domain.ErrRefreshTokenNotFound
	}
	return t, nil
}

func (r *fakeRefreshTokenRepo) Revoke(_ context.Context, tokenHash string) error {
	if t, ok := r.byHash[tokenHash]; ok {
		now := time.Now()
		t.RevokedAt = &now
	}
	return nil
}

func newTestAuthHandler() *AuthHandler {
	manager := jwt.NewManager("access-secret", "refresh-secret", 15*time.Minute, 168*time.Hour)
	authUsecase := usecases.NewAuthUsecase(newFakeUserRepo(), newFakeRefreshTokenRepo(), manager)
	return NewAuthHandler(authUsecase, false)
}

func newAuthTestRouter(h *AuthHandler) *gin.Engine {
	router := gin.New()
	router.POST("/auth/register", h.Register)
	router.POST("/auth/login", h.Login)
	router.POST("/auth/refresh", h.Refresh)
	router.POST("/auth/logout", h.Logout)
	return router
}

func doJSONRequest(router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func TestAuthHandler_Register_Success(t *testing.T) {
	router := newAuthTestRouter(newTestAuthHandler())

	w := doJSONRequest(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	})

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var body struct {
		Data authResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Data.User.Email != "user@example.com" {
		t.Errorf("user.email = %q, want %q", body.Data.User.Email, "user@example.com")
	}
	if body.Data.AccessToken == "" {
		t.Error("accessToken should not be empty")
	}

	refreshCookie := findCookie(w.Result().Cookies(), refreshTokenCookieName)
	if refreshCookie == nil {
		t.Fatal("expected a refresh_token cookie to be set")
	}
	if !refreshCookie.HttpOnly {
		t.Error("refresh_token cookie should be HttpOnly")
	}
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	router := newAuthTestRouter(newTestAuthHandler())

	w := doJSONRequest(router, http.MethodPost, "/auth/register", map[string]string{
		"password": "password123",
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}

	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Error.Code != string(apperrors.CodeValidation) {
		t.Errorf("error.code = %q, want %q", body.Error.Code, apperrors.CodeValidation)
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	router := newAuthTestRouter(newTestAuthHandler())

	doJSONRequest(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	})

	w := doJSONRequest(router, http.MethodPost, "/auth/login", map[string]string{
		"email":    "user@example.com",
		"password": "wrong-password",
	})

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusUnauthorized, w.Body.String())
	}
}

func TestAuthHandler_RefreshAndLogout(t *testing.T) {
	router := newAuthTestRouter(newTestAuthHandler())

	registerResp := doJSONRequest(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	})
	refreshCookie := findCookie(registerResp.Result().Cookies(), refreshTokenCookieName)
	if refreshCookie == nil {
		t.Fatal("expected a refresh_token cookie from register")
	}

	refreshReq := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	refreshReq.AddCookie(refreshCookie)
	refreshW := httptest.NewRecorder()
	router.ServeHTTP(refreshW, refreshReq)
	if refreshW.Code != http.StatusOK {
		t.Fatalf("refresh status = %d, want %d, body = %s", refreshW.Code, http.StatusOK, refreshW.Body.String())
	}

	rotatedCookie := findCookie(refreshW.Result().Cookies(), refreshTokenCookieName)
	if rotatedCookie == nil {
		t.Fatal("expected a rotated refresh_token cookie from refresh")
	}
	if rotatedCookie.Value == refreshCookie.Value {
		t.Error("refresh should rotate to a new token, not reuse the old one")
	}

	logoutReq := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	logoutReq.AddCookie(rotatedCookie)
	logoutW := httptest.NewRecorder()
	router.ServeHTTP(logoutW, logoutReq)
	if logoutW.Code != http.StatusOK {
		t.Fatalf("logout status = %d, want %d, body = %s", logoutW.Code, http.StatusOK, logoutW.Body.String())
	}

	reuseReq := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	reuseReq.AddCookie(rotatedCookie)
	reuseW := httptest.NewRecorder()
	router.ServeHTTP(reuseW, reuseReq)
	if reuseW.Code != http.StatusUnauthorized {
		t.Errorf("refresh after logout status = %d, want %d", reuseW.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Refresh_MissingCookie(t *testing.T) {
	router := newAuthTestRouter(newTestAuthHandler())

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
