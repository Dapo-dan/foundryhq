package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
)

const userIDContextKey = "userID"

const bearerPrefix = "Bearer "

// Auth validates the Authorization: Bearer <access_token> header using
// manager, rejecting the request with the standard unauthorized envelope
// (see docs/api.md) on anything missing, malformed, expired, or wrong
// secret/type. On success it sets the caller's user ID in context for
// downstream handlers/usecases to read via GetUserID. Routes that don't
// require auth (health checks, /auth/* itself) must not have this
// middleware applied — register them outside the group that uses it.
func Auth(manager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, bearerPrefix) || header == bearerPrefix {
			abortUnauthorized(c, "missing or malformed authorization header")
			return
		}

		claims, err := manager.ValidateAccessToken(strings.TrimPrefix(header, bearerPrefix))
		if err != nil {
			abortUnauthorized(c, "invalid or expired access token")
			return
		}

		c.Set(userIDContextKey, claims.UserID)
		c.Next()
	}
}

// GetUserID returns the user ID set by Auth. The bool is false if Auth
// hasn't run for this request — callers should treat that as a programmer
// error (a route registered without the middleware), not a 401.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get(userIDContextKey)
	if !exists {
		return uuid.UUID{}, false
	}
	id, ok := value.(uuid.UUID)
	return id, ok
}

func abortUnauthorized(c *gin.Context, message string) {
	appErr := apperrors.Unauthorized(message)
	c.AbortWithStatusJSON(apperrors.StatusFor(appErr.Code), gin.H{
		"error": gin.H{"code": appErr.Code, "message": appErr.Message},
	})
}
