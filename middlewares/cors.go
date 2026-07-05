package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const legacyAllowMethods = "GET, POST, PUT, DELETE, OPTIONS"

var legacyAllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Correlation-ID"}
var legacyExposeHeaders = []string{"X-Correlation-ID", "X-RateLimit-Remaining", "X-RateLimit-Limit"}

func CORSMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		applyLegacyCORSHeaders(c, origin)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func applyLegacyCORSHeaders(c *gin.Context, origin string) {
	// In production, use specific origins.
	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", legacyAllowMethods)
	c.Header("Access-Control-Allow-Headers", joinHeaderValues(legacyAllowHeaders))
	c.Header("Access-Control-Expose-Headers", joinHeaderValues(legacyExposeHeaders))
	c.Header("Access-Control-Allow-Credentials", "true")
}

func joinHeaderValues(values []string) string {
	if len(values) == 0 {
		return ""
	}

	joined := values[0]
	for i := 1; i < len(values); i++ {
		joined += ", " + values[i]
	}

	return joined
}

// Cors returns a CORS middleware configuration for the Gin router.
func Cors(allowedOrigins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
