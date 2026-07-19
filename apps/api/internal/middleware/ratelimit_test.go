package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// runRateLimit builds a one-route router with limiter's Limit applied and
// returns the recorded response for a single request from remoteAddr.
func runRateLimit(limiter *RateLimiter, remoteAddr string) *httptest.ResponseRecorder {
	router := gin.New()
	router.Use(limiter.Limit())
	router.POST("/login", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = remoteAddr
	router.ServeHTTP(w, req)
	return w
}

func TestRateLimiter_AllowsUpToBurst(t *testing.T) {
	limiter := NewRateLimiter(rate.Every(time.Minute), 3)

	for i := 0; i < 3; i++ {
		w := runRateLimit(limiter, "1.2.3.4:5678")
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d, want %d", i+1, w.Code, http.StatusOK)
		}
	}
}

func TestRateLimiter_RejectsOverBurst(t *testing.T) {
	limiter := NewRateLimiter(rate.Every(time.Minute), 3)

	for i := 0; i < 3; i++ {
		runRateLimit(limiter, "1.2.3.4:5678")
	}

	w := runRateLimit(limiter, "1.2.3.4:5678")
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimiter_TracksClientsIndependently(t *testing.T) {
	limiter := NewRateLimiter(rate.Every(time.Minute), 1)

	runRateLimit(limiter, "1.2.3.4:5678")
	w := runRateLimit(limiter, "5.6.7.8:5678")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d — a different client IP should have its own bucket", w.Code, http.StatusOK)
	}
}
