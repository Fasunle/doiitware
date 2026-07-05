// Package middlewares provides the Gin middleware implementations shared by Doiit services.
//
// The package focuses on the cross-cutting HTTP concerns that should behave the same way
// everywhere: authentication, recovery, request logging, request-size limits, security
// headers, input sanitization, tracing, and rate limiting.
//
// Common use cases include:
//
//   - protecting public APIs with JWT authentication and rate limiting
//   - standardizing browser security headers for front-end facing endpoints
//   - hiding internal stack traces while still logging enough context for operators
//
// Example:
//
//	router := gin.New()
//	limiter := helpers.NewIPRateLimiter(100, 200)
//
//	router.Use(middlewares.RequestLoggerMiddleware())
//	router.Use(middlewares.CustomRecoveryMiddleware())
//	router.Use(middlewares.RateLimitMiddleware(limiter, []string{"127.0.0.1"}))
//	router.Use(middlewares.SecurityHeadersMiddleware("https://api.example.com"))
package middlewares
