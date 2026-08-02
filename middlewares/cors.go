package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS is a middleware function that handles Cross-Origin Resource Sharing (CORS) for a Gin application.
// It allows you to specify which origins, methods, and headers are permitted for cross-origin requests.
// The middleware checks the request's origin against the allowed origins and sets the appropriate CORS headers in the response.
// If the origin is not allowed, it responds with a 403 Forbidden status. It also handles preflight OPTIONS requests by responding
// with a 204 No Content status.
//
// Allowed origins support wildcard patterns like:
//   - "https://*.example.com" matches all HTTPS subdomains of example.com
//   - "http://*.example.com" matches all HTTP subdomains of example.com
//   - "*.example.com" matches both HTTP and HTTPS subdomains
//
// Example usage:
//
//	router.Use(middlewares.CORS(
//		"https://example.com, https://api.example.com, https://*.example.com",
//		"GET, POST, PUT, DELETE",
//		nil,
//		nil,
//	))
func CORS(allowedOriginsStr string, methods string, headers *string, exposeHeaders *string) gin.HandlerFunc {
	// Pre-parse and validate origins
	allowedOrigins := ParseAllowedOrigins(allowedOriginsStr)

	if exposeHeaders == nil {
		exposeHeaders = new(string)
		*exposeHeaders = "Content-Length, X-Request-ID"
	}

	if headers == nil {
		headers = new(string)
		*headers = "Origin, Content-Type, Accept, Authorization, X-Request-ID"
	}

	// Ensure OPTIONS is included in the allowed methods for preflight requests
	if !strings.Contains(methods, "OPTIONS") {
		methods += ", OPTIONS"
	}

	return func(c *gin.Context) {
		reqOrigin := c.Request.Header.Get("Origin")

		// If no Origin header, skip CORS (it's a same-origin request)
		if reqOrigin == "" {
			c.Next()
			return
		}

		if !IsMethodAllowed(c.Request.Method, methods) {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}

		// Check if the request origin is allowed
		if !IsOriginAllowed(reqOrigin, allowedOrigins) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", reqOrigin)
		c.Header("Access-Control-Allow-Methods", methods)
		c.Header("Access-Control-Allow-Headers", *headers)
		c.Header("Access-Control-Expose-Headers", *exposeHeaders)
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ParseAllowedOrigins converts a comma-separated string into a list of origin patterns
func ParseAllowedOrigins(originsStr string) []string {
	if originsStr == "" {
		return []string{}
	}

	parts := strings.Split(originsStr, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

// IsOriginAllowed checks if an origin matches any allowed pattern
func IsOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if MatchOrigin(origin, allowed) {
			return true
		}
	}
	return false
}

// IsMethodAllowed checks if a method is in the allowed methods string
func IsMethodAllowed(method string, allowedMethods string) bool {
	method = strings.ToUpper(method)
	allowedMethods = strings.ToUpper(allowedMethods)
	return strings.Contains(allowedMethods, method)
}

// MatchOrigin handles wildcard matching with proper semantics
func MatchOrigin(origin, pattern string) bool {
	// Exact match
	if origin == pattern {
		return true
	}

	// Wildcard pattern: *.example.com or https://*.example.com
	if strings.Contains(pattern, "*") {
		// Extract protocol and domain from pattern
		patternProtocol := ""
		patternDomain := pattern
		if idx := strings.Index(pattern, "://"); idx != -1 {
			patternProtocol = pattern[:idx+3]
			patternDomain = pattern[idx+3:]
		}

		// Extract protocol and domain from origin
		originProtocol := ""
		originDomain := origin
		if idx := strings.Index(origin, "://"); idx != -1 {
			originProtocol = origin[:idx+3]
			originDomain = origin[idx+3:]
		}

		// If pattern has a protocol, it must match
		if patternProtocol != "" && patternProtocol != originProtocol {
			return false
		}

		// Handle wildcard domain patterns
		if strings.HasPrefix(patternDomain, "*.") {
			suffix := patternDomain[1:] // ".example.com"
			// Check if origin domain ends with the suffix
			if strings.HasSuffix(originDomain, suffix) {
				// Make sure it's not the naked domain
				subdomainPart := originDomain[:len(originDomain)-len(suffix)]
				// Must have at least one subdomain (not empty)
				if subdomainPart != "" {
					return true
				}
			}
		}
	}

	return false
}
