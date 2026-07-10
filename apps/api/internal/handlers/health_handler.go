package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler keeps its dependencies (here, db) as unexported fields set
// once at construction, rather than reaching for globals — this is the
// standard Go dependency-injection pattern and keeps the handler testable
// with a fake/mock db.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler constructs a HealthHandler. Go has no constructors, so
// a New<Type> factory function is the idiomatic stand-in.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check reports service health, including whether the database is reachable.
func (h *HealthHandler) Check(c *gin.Context) {
	status := "ok"
	code := http.StatusOK

	// gorm.DB doesn't expose a ping itself; DB() unwraps the underlying
	// database/sql handle, which is what actually knows how to check the
	// connection.
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
		status = "degraded"
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, gin.H{"status": status})
}
