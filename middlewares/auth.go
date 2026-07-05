package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fasunle/doiitware/cors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// AuthMiddleware validates Bearer tokens, loads the authenticated user into Gin context,
// and aborts the request with a structured 401 response when validation fails.
func AuthMiddleware(source string, public []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if source == "" {
			source = "service"
		}

		// check if the requested endpoint is a public route
		for _, route := range public {
			if strings.HasPrefix(c.Request.URL.Path, route) {
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "missing_authorization",
				"message": "Authorization header is required",
				"source":  source,
			})
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token_format",
				"message": "Invalid authorization header format. Use: Bearer <token>",
				"source":  source,
			})
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		claims := &cors.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "Invalid or expired token",
				"source":  source,
			})
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}
