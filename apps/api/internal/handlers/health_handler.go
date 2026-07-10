package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/version"
)

// HealthHandler serves the operational endpoints (/health, /ready, /version)
// that orchestrators, load balancers, and on-call engineers use to check
// whether the service is alive, ready for traffic, and which build is
// deployed. These are exempt from auth and API versioning (see docs/api.md).
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler constructs a HealthHandler. Go has no constructors, so
// a New<Type> factory function is the idiomatic stand-in.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health is a liveness check: it reports ok as long as the process is up
// and able to handle requests, with no dependency checks. An orchestrator
// should use this to decide whether to restart the container.
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Ready is a readiness check: it additionally verifies the database is
// reachable, since the API can't usefully serve traffic without it. A load
// balancer should use this to decide whether to route traffic here.
func (h *HealthHandler) Ready(c *gin.Context) {
	// gorm.DB doesn't expose a ping itself; DB() unwraps the underlying
	// database/sql handle, which is what actually knows how to check the
	// connection.
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Version reports the running build, useful for confirming a deploy landed
// or correlating a bug report with the exact commit in production.
func (h *HealthHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": version.Version,
		"commit":  version.Commit,
	})
}
