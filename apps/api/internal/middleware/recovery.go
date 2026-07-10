package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery catches panics from any handler and converts them into a 500
// using the API's standard error envelope (see docs/api.md's internal_error
// code), instead of gin's default recovery, which writes a bare stack trace
// to stderr and doesn't fit that envelope. The stack trace is logged
// server-side, tagged with the request ID for correlation.
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered",
					zap.String("request_id", GetRequestID(c)),
					zap.Any("panic", rec),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code":    "internal_error",
						"message": "internal server error",
					},
				})
			}
		}()

		c.Next()
	}
}
