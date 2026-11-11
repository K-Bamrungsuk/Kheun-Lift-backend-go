package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS lets the Next app (allowedOrigin) call this API with cookies.
func CORS(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Vary: Origin keeps caches from reusing one origin's CORS headers for another.
		c.Header("Vary", "Origin")
		if allowedOrigin != "" && c.GetHeader("Origin") == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Credentials", "true")
			if c.Request.Method == http.MethodOptions { // preflight
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
				c.Header("Access-Control-Allow-Headers", "Content-Type")
				c.Header("Access-Control-Max-Age", "600")
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}
		c.Next()
	}
}
