package middlewares

import (
	"fmt"

	"github.com/fasunle/doiitware/config"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds the shared browser security headers to every response.
//
// The route value is used in the Content-Security-Policy connect-src directive.
func SecurityHeadersMiddleware(route string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// HSTS - Force HTTPS (preload recommended)
		c.Header(config.HeaderStrictTransport, "max-age=31536000; includeSubDomains; preload")

		// CSP - Mitigate XSS and data injection
		c.Header(config.HeaderContentSecurity, "default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; "+
			"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; "+
			"img-src 'self' data: https:; "+
			"font-src 'self' https://cdn.jsdelivr.net; "+
			fmt.Sprintf("connect-src 'self' %s; ", route)+
			"frame-ancestors 'none'; "+
			"form-action 'self'")

		// X-Frame-Options - Prevent clickjacking
		c.Header(config.HeaderXFrameOptions, "DENY")

		// X-Content-Type-Options - Prevent MIME sniffing
		c.Header(config.HeaderXContentType, "nosniff")

		// Referrer-Policy - Control referrer information
		c.Header(config.HeaderReferrerPolicy, "strict-origin-when-cross-origin")

		// Permissions-Policy - Restrict browser features
		c.Header(config.HeaderPermissionsPolicy, "geolocation=(), microphone=(), camera=(), payment=()")

		// X-XSS-Protection - Legacy XSS protection
		c.Header(config.HeaderXSSProtection, "1; mode=block")

		// Remove Server header to hide technology
		c.Header("Server", "")

		c.Next()
	}
}
