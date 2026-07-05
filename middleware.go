package doiitmiddleware

import (
	"strings"

	"github.com/fasunle/doiitware/helpers"
	middleware "github.com/fasunle/doiitware/middlewares"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var DefaultRateLimitWhitelist = []string{
	"127.0.0.1",
	"::1",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"localhost",
}

func NewIPRateLimiter(limitsPerSecond rate.Limit, burst int) *helpers.IPRateLimiter {
	return helpers.NewIPRateLimiter(limitsPerSecond, burst)
}

func RateLimitMiddleware(rateLimiter *helpers.IPRateLimiter, whitelistCIDRs []string) gin.HandlerFunc {
	if len(whitelistCIDRs) == 0 {
		whitelistCIDRs = DefaultRateLimitWhitelist
	}

	return middleware.RateLimitMiddleware(rateLimiter, whitelistCIDRs)
}

func CustomRecoveryMiddleware() gin.HandlerFunc {
	return middleware.CustomRecoveryMiddleware()
}

func SecurityHeadersMiddleware(route string) gin.HandlerFunc {
	return middleware.SecurityHeadersMiddleware(route)
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return middleware.RequestLoggerMiddleware()
}

func RequestSizeMiddleware() gin.HandlerFunc {
	return middleware.RequestSizeMiddleware()
}

func SanitizeInputMiddleware() gin.HandlerFunc {
	return middleware.SanitizeInputMiddleware()
}

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
