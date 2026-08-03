package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders sets a baseline of response headers that cost nothing for
// a JSON API to send but close off MIME-sniffing and clickjacking vectors
// against anything that ends up fronting or embedding it later. hsts should
// be true only when the API is actually served over HTTPS (production) —
// sending Strict-Transport-Security over plain HTTP local dev would be
// meaningless and could confuse a browser that later hits the same host.
func SecurityHeaders(hsts bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		if hsts {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		c.Next()
	}
}
