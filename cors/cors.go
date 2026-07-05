package cors

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/fasunle/doiitware/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	ClientID string   `json:"client_id,omitempty"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// ProtectedHandler requires authentication
func ProtectedHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	email := c.GetString("email")
	roles := c.GetStringSlice("roles")

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the protected endpoint!",
		"user": gin.H{
			"id":    userID,
			"email": email,
			"roles": roles,
		},
	})
}

// GenerateJWT creates a new JWT token
func GenerateJWT(userID, email string, roles []string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.JwtIssuer,
			Audience:  []string{config.JwtAudience},
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JwtExpiry)),
			ID:        generateTokenID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}

// generateTokenID creates a unique token ID
func generateTokenID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
