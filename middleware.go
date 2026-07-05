package doiitware

import (
	"context"
	"strings"

	"github.com/fasunle/doiitware/helpers"
	middleware "github.com/fasunle/doiitware/middlewares"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// DefaultRateLimitWhitelist contains CIDR ranges and hostnames that bypass the rate limiter.
//
// It is intended for loopback and private-network traffic in local development and
// internal service-to-service calls.
var DefaultRateLimitWhitelist = []string{
	"127.0.0.1",
	"::1",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"localhost",
}

// NewIPRateLimiter creates an IP-scoped token bucket limiter with the provided rate and burst.
func NewIPRateLimiter(limitsPerSecond rate.Limit, burst int) *helpers.IPRateLimiter {
	return helpers.NewIPRateLimiter(limitsPerSecond, burst)
}

// RateLimitMiddleware wraps the shared rate limiter middleware and applies the default whitelist
// when no CIDRs are provided.
func RateLimitMiddleware(rateLimiter *helpers.IPRateLimiter, whitelistCIDRs []string) gin.HandlerFunc {
	if len(whitelistCIDRs) == 0 {
		whitelistCIDRs = DefaultRateLimitWhitelist
	}

	return middleware.RateLimitMiddleware(rateLimiter, whitelistCIDRs)
}

// CustomRecoveryMiddleware returns a panic recovery middleware that hides internal errors from clients.
func CustomRecoveryMiddleware() gin.HandlerFunc {
	return middleware.CustomRecoveryMiddleware()
}

// SecurityHeadersMiddleware adds the shared browser security headers for the given route or origin.
func SecurityHeadersMiddleware(route string) gin.HandlerFunc {
	return middleware.SecurityHeadersMiddleware(route)
}

// RequestLoggerMiddleware returns the shared request logging middleware.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return middleware.RequestLoggerMiddleware()
}

// RequestSizeMiddleware limits request bodies to the shared maximum size.
func RequestSizeMiddleware() gin.HandlerFunc {
	return middleware.RequestSizeMiddleware()
}

// SanitizeInputMiddleware returns the shared input sanitization middleware.
func SanitizeInputMiddleware() gin.HandlerFunc {
	return middleware.SanitizeInputMiddleware()
}

// TracingMiddleware returns the shared tracing middleware, which can be enabled or disabled and accepts a custom tracer function.
func TracingMiddleware(enabled bool, tracer func(context.Context) error) gin.HandlerFunc {
	return middleware.TracingMiddleware(enabled, tracer)
}

// ParseCIDRs splits a comma-separated CIDR list into a normalized slice.
func ParseCIDRs(cidrs string) []string {
	parts := strings.Split(cidrs, ",")
	normalized := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}

	return normalized
}
