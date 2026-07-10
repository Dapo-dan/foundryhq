// Package middleware holds cross-cutting Gin middleware (auth, CORS,
// logging, request ID). It must not contain business logic — see
// docs/architecture.md.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHeader is the header carrying the request ID, both inbound (an
// upstream proxy may already have set one) and outbound (so the caller can
// correlate their request with server-side logs).
const RequestIDHeader = "X-Request-ID"

const requestIDContextKey = "requestID"

// RequestID assigns a unique ID to each request, reusing an inbound
// X-Request-ID if one is already present rather than overwriting it — that
// lets a request keep the same ID as it crosses multiple services. It must
// run before Logger so the ID is available to include in request logs.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}

		c.Set(requestIDContextKey, id)
		c.Writer.Header().Set(RequestIDHeader, id)
		c.Next()
	}
}

// GetRequestID returns the request ID set by RequestID, or "" if RequestID
// hasn't run (e.g. called outside a request, or before the middleware).
func GetRequestID(c *gin.Context) string {
	id, _ := c.Get(requestIDContextKey)
	s, _ := id.(string)
	return s
}
