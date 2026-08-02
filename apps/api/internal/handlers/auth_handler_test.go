package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func (r *fakeUserRepo) UpdatePassword(_ context.Context, userID uuid.UUID, passwordHash string) error {
	u, ok := r.byID[userID]
	if !ok {
		return domain.ErrUserNotFound
	}
	u.PasswordHash = passwordHash
	return nil
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

func (r *fakeRefreshTokenRepo) RevokeAllForUser(_ context.Context, userID uuid.UUID) error {
	now := time.Now()
	for _, t := range r.byHash {
		if t.UserID == userID && t.RevokedAt == nil {
			t.RevokedAt = &now
		}
	}
	return nil
}

type fakePasswordResetTokenRepo struct {
	byHash map[string]*domain.PasswordResetToken
}

func newFakePasswordResetTokenRepo() *fakePasswordResetTokenRepo {
	return &fakePasswordResetTokenRepo{byHash: map[string]*domain.PasswordResetToken{}}
}

func (r *fakePasswordResetTokenRepo) Create(_ context.Context, token *domain.PasswordResetToken) error {
	token.ID = uuid.New()
	r.byHash[token.TokenHash] = token
	return nil
}

func (r *fakePasswordResetTokenRepo) GetByTokenHash(_ context.Context, tokenHash string) (*domain.PasswordResetToken, error) {
	t, ok := r.byHash[tokenHash]
	if !ok {
		return nil, domain.ErrPasswordResetTokenNotFound
	}
	return t, nil
}

func (r *fakePasswordResetTokenRepo) MarkUsed(_ context.Context, tokenHash string) error {
	if t, ok := r.byHash[tokenHash]; ok {
		now := time.Now()
		t.UsedAt = &now
	}
	return nil
}

// fakeMailer is a mailer.EmailSender double that records every send so
// tests can pull the emailed reset token out of it instead of hitting a
// real provider.
type fakeMailer struct {
	sent []fakeSentEmail
}

type fakeSentEmail struct {
	to, subject, body string
}

func (m *fakeMailer) Send(_ context.Context, to, subject, body string) error {
	m.sent = append(m.sent, fakeSentEmail{to: to, subject: subject, body: body})
	return nil
}

func newTestAuthHandler() *AuthHandler {
	h, _ := newTestAuthHandlerWithMailer()
	return h
}

func newTestAuthHandlerWithMailer() (*AuthHandler, *fakeMailer) {
	h, mailer, _, _, _ := newTestAuthHandlerFull()
	return h, mailer
}

// newTestAuthHandlerFull is the one place every AuthHandler test dependency
// gets constructed — other helpers above just narrow which of these they
// return.
func newTestAuthHandlerFull() (*AuthHandler, *fakeMailer, *fakeInviteTokenRepo, *fakeWorkspaceMemberRepo, *fakeUserRepo) {
	manager := jwt.NewManager("access-secret", "refresh-secret", 15*time.Minute, 168*time.Hour)
	mailer := &fakeMailer{}
	inviteTokens := newFakeInviteTokenRepo()
	members := newFakeWorkspaceMemberRepo()
	users := newFakeUserRepo()
	authUsecase := usecases.NewAuthUsecase(
		users,
		newFakeRefreshTokenRepo(),
		newFakePasswordResetTokenRepo(),
		inviteTokens,
		members,
		manager,
		mailer,
		"http://localhost:5173",
	)
	return NewAuthHandler(authUsecase, false), mailer, inviteTokens, members, users
}

func newAuthTestRouter(h *AuthHandler) *gin.Engine {
	router := gin.New()
	router.POST("/auth/register", h.Register)
	router.POST("/auth/login", h.Login)
	router.POST("/auth/refresh", h.Refresh)
	router.POST("/auth/logout", h.Logout)
	router.POST("/auth/forgot-password", h.ForgotPassword)
	router.POST("/auth/reset-password", h.ResetPassword)
	router.POST("/auth/accept-invite", h.AcceptInvite)
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

func TestAuthHandler_ForgotPassword_AlwaysOK(t *testing.T) {
	h, mailer := newTestAuthHandlerWithMailer()
	router := newAuthTestRouter(h)

	doJSONRequest(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	})

	registered := doJSONRequest(router, http.MethodPost, "/auth/forgot-password", map[string]string{
		"email": "user@example.com",
	})
	if registered.Code != http.StatusOK {
		t.Fatalf("registered-email status = %d, want %d, body = %s", registered.Code, http.StatusOK, registered.Body.String())
	}

	unregistered := doJSONRequest(router, http.MethodPost, "/auth/forgot-password", map[string]string{
		"email": "nobody@example.com",
	})
	if unregistered.Code != http.StatusOK {
		t.Fatalf("unregistered-email status = %d, want %d (enumeration-safe)", unregistered.Code, http.StatusOK)
	}

	if len(mailer.sent) != 1 {
		t.Fatalf("expected exactly 1 sent email, got %d", len(mailer.sent))
	}
}

func TestAuthHandler_ResetPassword_Success(t *testing.T) {
	h, mailer := newTestAuthHandlerWithMailer()
	router := newAuthTestRouter(h)

	doJSONRequest(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	})
	doJSONRequest(router, http.MethodPost, "/auth/forgot-password", map[string]string{
		"email": "user@example.com",
	})
	if len(mailer.sent) != 1 {
		t.Fatalf("expected a reset email to have been sent, got %d", len(mailer.sent))
	}
	token := extractToken(t, mailer.sent[0].body)

	w := doJSONRequest(router, http.MethodPost, "/auth/reset-password", map[string]string{
		"token":    token,
		"password": "new-password456",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	loginW := doJSONRequest(router, http.MethodPost, "/auth/login", map[string]string{
		"email":    "user@example.com",
		"password": "new-password456",
	})
	if loginW.Code != http.StatusOK {
		t.Errorf("login with new password status = %d, want %d", loginW.Code, http.StatusOK)
	}
}

func TestAuthHandler_AcceptInvite_Success(t *testing.T) {
	h, _, inviteTokens, members, users := newTestAuthHandlerFull()
	router := newAuthTestRouter(h)

	placeholder := &domain.User{Email: "invited@example.com"}
	if err := users.Create(context.Background(), placeholder); err != nil {
		t.Fatalf("seeding placeholder user: %v", err)
	}
	member := &domain.WorkspaceMember{WorkspaceID: uuid.New(), UserID: placeholder.ID, Role: domain.RoleMember}
	if err := members.Create(context.Background(), member); err != nil {
		t.Fatalf("seeding placeholder membership: %v", err)
	}
	rawToken := "raw-invite-token"
	if err := inviteTokens.Create(context.Background(), &domain.InviteToken{
		WorkspaceMemberID: member.ID,
		TokenHash:         sha256Hex(rawToken),
		ExpiresAt:         time.Now().Add(time.Hour),
	}); err != nil {
		t.Fatalf("seeding invite token: %v", err)
	}

	w := doJSONRequest(router, http.MethodPost, "/auth/accept-invite", map[string]string{
		"token":    rawToken,
		"password": "password123",
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	refreshCookie := findCookie(w.Result().Cookies(), refreshTokenCookieName)
	if refreshCookie == nil {
		t.Fatal("expected a refresh_token cookie to be set")
	}
}

func TestAuthHandler_AcceptInvite_InvalidToken(t *testing.T) {
	router := newAuthTestRouter(newTestAuthHandler())

	w := doJSONRequest(router, http.MethodPost, "/auth/accept-invite", map[string]string{
		"token":    "not-a-real-token",
		"password": "password123",
	})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusUnauthorized, w.Body.String())
	}
}

func TestAuthHandler_ResetPassword_InvalidToken(t *testing.T) {
	router := newAuthTestRouter(newTestAuthHandler())

	w := doJSONRequest(router, http.MethodPost, "/auth/reset-password", map[string]string{
		"token":    "not-a-real-token",
		"password": "new-password456",
	})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusUnauthorized, w.Body.String())
	}
}

// sha256Hex mirrors usecases.hashToken's (unexported, different package)
// algorithm so handler tests can seed a token whose hash the real usecase
// will recognize — usecases never persists or exposes the raw token itself.
func sha256Hex(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// extractToken pulls the raw reset token out of the query string of the
// link embedded in a ForgotPassword email body.
func extractToken(t *testing.T, body string) string {
	t.Helper()
	const marker = "token="
	i := strings.Index(body, marker)
	if i == -1 {
		t.Fatalf("email body has no %q: %s", marker, body)
	}
	rest := body[i+len(marker):]
	end := strings.IndexAny(rest, `"'&`)
	if end == -1 {
		end = len(rest)
	}
	return rest[:end]
}
