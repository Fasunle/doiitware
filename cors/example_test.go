package cors

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ExampleGenerateJWT() {
	jwtSecret = []byte("example-secret")
	_, _ = GenerateJWT("user-123", "user@example.com", []string{"admin"})
}

func ExampleProtectedHandler() {
	router := gin.New()
	router.GET("/me", ProtectedHandler)
	_ = router
}

func ExampleClaims() {
	_ = Claims{UserID: "user-123", Email: "user@example.com", Roles: []string{"admin"}, RegisteredClaims: jwt.RegisteredClaims{}}
}
