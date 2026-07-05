package middlewares

import (
	"log"
	"time"

	"github.com/fasunle/doiitware/helpers"
	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware attaches a request ID, then logs the request method, path,
// status, latency, client IP, and user agent after the handler completes.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID for tracing
		requestID := helpers.GenerateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		start := time.Now()

		c.Next()

		// Log after request completes
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %s | Status: %d | Latency: %v | IP: %s | User-Agent: %s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			latency,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}
