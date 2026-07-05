package middlewares

import (
	"context"
	"log"

	"github.com/fasunle/doiitware/helpers"
	"github.com/gin-gonic/gin"
)

func TracingMiddleware(enabled bool, tracer func(context.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate correlation ID
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = helpers.GenerateCorrelationID()
		}

		c.Set("correlation_id", correlationID)
		c.Header("X-Correlation-ID", correlationID)
		c.Request.Header.Set("X-Correlation-ID", correlationID)

		// Add span context if tracing enabled
		if enabled {
			// Start span
			ctx := context.Background()
			if err := tracer(ctx); err != nil {
				log.Panicf("Failed to set up tracer: %v", err)
			}
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}
