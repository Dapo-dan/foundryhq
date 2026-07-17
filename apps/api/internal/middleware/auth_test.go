package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func testManager() *jwt.Manager {
	return jwt.NewManager("access-secret", "refresh-secret", 15*time.Minute, 168*time.Hour)
}

// runAuth builds a one-route router with Auth applied and returns the
// recorded response plus the userID (if any) the downstream handler saw.
func runAuth(manager *jwt.Manager, authHeader string) (w *httptest.ResponseRecorder, sawUserID uuid.UUID, sawOK bool) {
	router := gin.New()
	router.Use(Auth(manager))
	router.GET("/protected", func(c *gin.Context) {
		sawUserID, sawOK = GetUserID(c)
		c.Status(http.StatusOK)
	})

	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	router.ServeHTTP(w, req)
	return w, sawUserID, sawOK
}

func TestAuth_ValidToken(t *testing.T) {
	manager := testManager()
	userID := uuid.New()
	token, err := manager.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	w, sawUserID, sawOK := runAuth(manager, bearerPrefix+token)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if !sawOK {
		t.Fatal("GetUserID() ok = false, want true")
	}
	if sawUserID != userID {
		t.Errorf("GetUserID() = %v, want %v", sawUserID, userID)
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	w, _, sawOK := runAuth(testManager(), "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if sawOK {
		t.Error("downstream handler should not have run")
	}
}

func TestAuth_MalformedHeader(t *testing.T) {
	w, _, _ := runAuth(testManager(), "Token abc123")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	manager := jwt.NewManager("access-secret", "refresh-secret", -1*time.Minute, 168*time.Hour)
	token, err := manager.GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	w, _, _ := runAuth(manager, bearerPrefix+token)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_RefreshTokenRejected(t *testing.T) {
	manager := testManager()
	refreshToken, _, err := manager.GenerateRefreshToken(uuid.New())
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	w, _, _ := runAuth(manager, bearerPrefix+refreshToken)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if _, ok := GetUserID(c); ok {
		t.Error("GetUserID() ok = true, want false when Auth hasn't run")
	}
}
