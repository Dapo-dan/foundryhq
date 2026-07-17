package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
)

// handleError maps any error returned by a usecase to the standard error
// envelope ({"error":{"code","message","field?"}}, see docs/api.md). Errors
// that aren't *apperrors.Error are treated as unexpected: they're recorded
// on the gin context (surfaced by middleware.Logger with the request ID)
// and never serialized beyond the generic internal_error message, so raw
// GORM/SQL errors never leak to clients.
func handleError(c *gin.Context, err error) {
	var appErr *apperrors.Error
	if errors.As(err, &appErr) {
		body := gin.H{"code": appErr.Code, "message": appErr.Message}
		if appErr.Field != "" {
			body["field"] = appErr.Field
		}
		c.JSON(apperrors.StatusFor(appErr.Code), gin.H{"error": body})
		return
	}

	c.Error(err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{"code": apperrors.CodeInternal, "message": "internal server error"},
	})
}
