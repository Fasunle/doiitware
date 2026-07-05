package middlewares

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fasunle/doiitware/helpers"
)

// RateLimitMiddleware enforces per-client request limits and bypasses whitelisted CIDRs.
func RateLimitMiddleware(rateLimiter *helpers.IPRateLimiter, whitelistCIDRs []string) gin.HandlerFunc {
	allowedCIDRs := normalizeCIDRs(whitelistCIDRs)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Check whitelist (for internal services)
		for _, cidr := range allowedCIDRs {
			if ipInCIDR(ip, strings.TrimSpace(cidr)) {
				c.Next()
				return
			}
		}

		clientLimiter := rateLimiter.GetLimiter(ip)
		if !clientLimiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate_limit_exceeded",
				"message":     "Too many requests. Please slow down.",
				"retry_after": time.Now().Add(time.Second).Unix(),
			})
			return
		}
		c.Next()
	}
}

func normalizeCIDRs(cidrs []string) []string {
	normalized := make([]string, 0, len(cidrs))

	for _, cidr := range cidrs {
		trimmed := strings.TrimSpace(cidr)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}

	return normalized
}

// ipInCIDR reports whether the provided IP belongs to the given CIDR range.
func ipInCIDR(ipStr, cidr string) bool {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(ipStr)
	return ipNet.Contains(ip)
}
