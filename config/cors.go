package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors returns a CORS middleware configuration for the Gin router.
func Cors(PORT string, allowedOrigins []string) gin.HandlerFunc {
	if allowedOrigins == nil {
		allowedOrigins = []string{}

		if os.Getenv("ENVIRONMENT") == "development" {
			allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://localhost%s", PORT))
		} else if os.Getenv("ENVIRONMENT") == "production" {
			allowedOrigins = append(allowedOrigins, "https://doiit.com")
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
