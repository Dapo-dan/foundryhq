package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
)

// visitorIdleTimeout bounds how long a client IP's bucket is kept after its
// last request, so the visitors map doesn't grow forever over the server's
// lifetime.
const visitorIdleTimeout = 10 * time.Minute

// visitor pairs a client's token bucket with when it was last seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter grants each client IP its own token bucket, used to slow down
// brute-force attempts against unauthenticated endpoints like /auth/login.
// It's in-memory and per-process — fine for the current single-instance
// deployment; a shared store (Redis) would be needed if the API ever runs
// more than one replica, since each replica would otherwise track its own
// independent buckets.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    rate.Limit
	burst    int
}

// NewRateLimiter builds a RateLimiter allowing burst requests immediately
// per client IP, then refilling at limit tokens/sec. It starts a background
// goroutine that evicts idle visitors so memory use stays bounded; the
// goroutine runs for the lifetime of the process, matching the server's own
// lifetime.
func NewRateLimiter(limit rate.Limit, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		burst:    burst,
	}
	go rl.evictIdleVisitors()
	return rl
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[key]
	if !ok {
		v = &visitor{limiter: rate.NewLimiter(rl.limit, rl.burst)}
		rl.visitors[key] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) evictIdleVisitors() {
	for range time.Tick(time.Minute) {
		rl.mu.Lock()
		for key, v := range rl.visitors {
			if time.Since(v.lastSeen) > visitorIdleTimeout {
				delete(rl.visitors, key)
			}
		}
		rl.mu.Unlock()
	}
}

// Limit returns middleware that rejects a request with 429 once the
// caller's client IP has exhausted its token bucket. It's rejected before
// reaching the handler, so a flood of attempts never touches the usecase
// (no password hashing, no DB query) — the whole point is to make
// brute-forcing expensive at the network edge, not just slow downstream.
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.getLimiter(c.ClientIP()).Allow() {
			appErr := apperrors.RateLimited("too many attempts, please try again later")
			c.AbortWithStatusJSON(apperrors.StatusFor(appErr.Code), gin.H{
				"error": gin.H{"code": appErr.Code, "message": appErr.Message},
			})
			return
		}
		c.Next()
	}
}
