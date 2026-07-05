package doiitware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func Example() {
	router := gin.New()
	limiter := NewIPRateLimiter(rate.Limit(100), 200)

	router.Use(RequestLoggerMiddleware())
	router.Use(CustomRecoveryMiddleware())
	router.Use(RateLimitMiddleware(limiter, nil))
	router.Use(SecurityHeadersMiddleware("https://api.example.com"))
	router.Use(TracingMiddleware(true, nil))

	_ = router
}

func ExampleParseCIDRs() {
	_ = ParseCIDRs("127.0.0.1, ::1, 10.0.0.0/8")
}