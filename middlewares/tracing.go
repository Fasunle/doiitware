package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware logs all requests with performance metrics
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID for tracing
		requestID := generateRequestID()
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

// generateRequestID creates a unique request ID
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
