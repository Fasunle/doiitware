package middlewares

import (
	"net/http"
	"strings"

	"github.com/fasunle/doiitware/config"

	"github.com/gin-gonic/gin"
)

// SanitizeInputMiddleware sanitizes query parameters and path parameters before handlers run.
func SanitizeInputMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		query := c.Request.URL.Query()
		for key, values := range query {
			if len(values) > 0 {
				query.Set(key, sanitizeString(values[0]))
			}
		}
		c.Request.URL.RawQuery = query.Encode()

		// Sanitize path parameters
		for _, param := range c.Params {
			param.Value = sanitizeString(param.Value)
		}

		c.Next()
	}
}

// sanitizeString removes characters commonly used in trivial HTML and injection payloads.
func sanitizeString(s string) string {
	// Basic sanitization - in production, use a proper library
	// Remove script tags, HTML tags, and common injection patterns
	replacements := map[string]string{
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
		"'":  "&#39;",
		"&":  "&amp;",
		";":  "",
		"/*": "",
		"*/": "",
		"--": "",
	}

	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	return s
}

// RequestSizeMiddleware caps the request body size at the shared platform limit.
func RequestSizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxRequestSize)
		c.Next()
	}
}
