package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// CustomRecoveryMiddleware recovers from panics, logs request context, and returns a safe error.
func CustomRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := string(debug.Stack())

				// Log the error with context
				log.Printf("PANIC RECOVERED: %v\nRequest: %s %s\nIP: %s\nStack: %s\n",
					err,
					c.Request.Method,
					c.Request.URL.Path,
					c.ClientIP(),
					stack,
				)

				// Send error response (don't expose internals)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":      "internal_server_error",
					"message":    "An unexpected error occurred. Our team has been notified.",
					"request_id": c.GetString("request_id"),
				})

				// Log to monitoring system (e.g., Sentry)
				// sentry.CaptureException(fmt.Errorf("%v", err))
			}
		}()
		c.Next()
	}
}
