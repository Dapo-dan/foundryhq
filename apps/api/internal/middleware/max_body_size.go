package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxBodySize caps every request body at limit bytes, regardless of what a
// client claims or actually sends. Without it, a caller — including an
// unauthenticated one, since /auth/register has no auth gate in front of
// it — can make the server buffer an arbitrarily large body in memory
// before request binding even gets a chance to reject it as invalid.
func MaxBodySize(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}
