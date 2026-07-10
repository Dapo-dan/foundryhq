package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS restricts cross-origin requests to allowedOrigins (the web/mobile
// client origins). AllowCredentials is required for the httpOnly
// refresh-token cookie (see adr/0004-jwt-access-refresh-tokens.md), which in
// turn means an explicit origin list must be used — the CORS spec forbids
// combining AllowCredentials with a "*" wildcard origin.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", RequestIDHeader},
		ExposeHeaders:    []string{RequestIDHeader},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
