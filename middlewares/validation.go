package middleware

import (
	"net/http"
	"strings"

	"github.com/fasunle/doiitware/config"

	"github.com/gin-gonic/gin"
)

// SanitizeInputMiddleware sanitizes request inputs to prevent XSS and injection
func SanitizeInputMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				c.Request.URL.Query().Set(key, sanitizeString(values[0]))
			}
		}

		// Sanitize path parameters
		for _, param := range c.Params {
			param.Value = sanitizeString(param.Value)
		}

		c.Next()
	}
}

// sanitizeString removes potentially dangerous characters
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

// RequestSizeMiddleware limits request body size
func RequestSizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxRequestSize)
		c.Next()
	}
}
