package middlewares

import (
	"github.com/fasunle/doiitware/helpers"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func ExampleRequestLoggerMiddleware() {
	router := gin.New()
	router.Use(RequestLoggerMiddleware())
	_ = router
}

func ExampleCustomRecoveryMiddleware() {
	router := gin.New()
	router.Use(CustomRecoveryMiddleware())
	_ = router
}

func ExampleRateLimitMiddleware() {
	router := gin.New()
	limiter := helpers.NewIPRateLimiter(rate.Limit(50), 100)
	router.Use(RateLimitMiddleware(limiter, []string{"127.0.0.1"}))
	_ = router
}
