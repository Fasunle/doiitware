// Package doiitware provides shared middleware, helper primitives, and constants for Doiit services.
//
// The package is designed to keep API gateway and service binaries aligned on core
// HTTP behavior such as rate limiting, request logging, panic recovery, security
// headers, request-size enforcement, and input sanitization.
//
// Typical use cases include:
//
//   - wrapping a Gin router with the same middleware stack in every service
//   - reusing the shared IP rate limiter and whitelist defaults
//   - applying consistent request IDs and recovery behavior across public APIs
//
// Example:
//
//   package main
//
//   import (
//       doiitware "github.com/fasunle/doiitware"
//       "github.com/gin-gonic/gin"
//   )
//
//   func main() {
//       router := gin.New()
//       limiter := doiitware.NewIPRateLimiter(100, 200)
//
//       router.Use(doiitware.RequestLoggerMiddleware())
//       router.Use(doiitware.CustomRecoveryMiddleware())
//       router.Use(doiitware.RateLimitMiddleware(limiter, nil))
//       router.Use(doiitware.SecurityHeadersMiddleware("https://api.example.com"))
//
//       _ = router.Run()
//   }
package doiitware
