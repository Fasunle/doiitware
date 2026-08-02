package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fasunle/doiitware/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		allowedOriginsStr string
		methods           string
		headers           *string
		exposeHeaders     *string
		requestOrigin     string
		requestMethod     string
		expectedStatus    int
		expectedHeaders   map[string]string
	}{
		{
			name:              "allowed exact origin - GET request",
			allowedOriginsStr: "https://example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Expose-Headers":    "Content-Length, X-Request-ID",
			},
		},
		{
			name:              "allowed exact origin - OPTIONS preflight",
			allowedOriginsStr: "https://example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://example.com",
			requestMethod:     "OPTIONS",
			expectedStatus:    http.StatusNoContent,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "https://example.com",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			},
		},
		{
			name:              "disallowed origin",
			allowedOriginsStr: "https://example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://evil.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusForbidden,
			expectedHeaders:   map[string]string{},
		},
		{
			name:              "wildcard subdomain match - simple subdomain",
			allowedOriginsStr: "https://*.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://api.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://api.example.com",
				"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:              "wildcard subdomain match - deep subdomain",
			allowedOriginsStr: "https://*.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://v1.api.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://v1.api.example.com",
				"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:              "wildcard subdomain match - different subdomain",
			allowedOriginsStr: "https://*.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://dashboard.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "https://dashboard.example.com",
			},
		},
		{
			name:              "wildcard subdomain no match - naked domain",
			allowedOriginsStr: "https://*.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusForbidden,
			expectedHeaders:   map[string]string{},
		},
		{
			name:              "wildcard with HTTP protocol",
			allowedOriginsStr: "http://*.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "http://api.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "http://api.example.com",
			},
		},
		{
			name:              "wildcard protocol mismatch",
			allowedOriginsStr: "https://*.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "http://api.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusForbidden,
			expectedHeaders:   map[string]string{},
		},
		{
			name:              "wildcard without protocol",
			allowedOriginsStr: "*.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://api.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "https://api.example.com",
			},
		},
		{
			name:              "wildcard without protocol with HTTP",
			allowedOriginsStr: "*.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "http://api.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "http://api.example.com",
			},
		},
		{
			name:              "multiple origins - first allowed",
			allowedOriginsStr: "https://app.example.com, https://admin.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://app.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "https://app.example.com",
			},
		},
		{
			name:              "multiple origins - second allowed",
			allowedOriginsStr: "https://app.example.com, https://admin.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://admin.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "https://admin.example.com",
			},
		},
		{
			name:              "mixed wildcard and exact",
			allowedOriginsStr: "https://app.example.com, https://*.api.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://v1.api.example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "https://v1.api.example.com",
			},
		},
		{
			name:              "mixed wildcard and exact - no match",
			allowedOriginsStr: "https://app.example.com, https://*.api.example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://evil.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusForbidden,
			expectedHeaders:   map[string]string{},
		},
		{
			name:              "no origin header - same origin request",
			allowedOriginsStr: "https://example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders:   map[string]string{},
		},
		{
			name:              "custom headers and expose headers",
			allowedOriginsStr: "https://example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           stringPtr("X-Custom-Header, Authorization"),
			exposeHeaders:     stringPtr("X-Custom-Response, X-Request-ID"),
			requestOrigin:     "https://example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Allow-Headers":     "X-Custom-Header, Authorization",
				"Access-Control-Expose-Headers":    "X-Custom-Response, X-Request-ID",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:              "origin with trailing slash",
			allowedOriginsStr: "https://example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "https://example.com/",
			requestMethod:     "GET",
			expectedStatus:    http.StatusForbidden,
			expectedHeaders:   map[string]string{},
		},
		{
			name:              "HTTP vs HTTPS mismatch",
			allowedOriginsStr: "https://example.com",
			methods:           "GET, POST, PUT, DELETE",
			headers:           nil,
			exposeHeaders:     nil,
			requestOrigin:     "http://example.com",
			requestMethod:     "GET",
			expectedStatus:    http.StatusForbidden,
			expectedHeaders:   map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middlewares.CORS(tt.allowedOriginsStr, tt.methods, tt.headers, tt.exposeHeaders))

			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			router.OPTIONS("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.requestMethod, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			for key, expectedValue := range tt.expectedHeaders {
				assert.Equal(t, expectedValue, w.Header().Get(key), "Header %s mismatch", key)
			}
		})
	}
}

func TestIsMethodAllowed(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		allowedMethods string
		expected       bool
	}{
		{
			name:           "GET is allowed",
			method:         "GET",
			allowedMethods: "GET, POST, PUT, DELETE",
			expected:       true,
		},
		{
			name:           "POST is allowed",
			method:         "POST",
			allowedMethods: "GET, POST, PUT, DELETE",
			expected:       true,
		},
		{
			name:           "PATCH is not allowed",
			method:         "PATCH",
			allowedMethods: "GET, POST, PUT, DELETE",
			expected:       false,
		},
		{
			name:           "case insensitive",
			method:         "get",
			allowedMethods: "GET, POST, PUT, DELETE",
			expected:       true,
		},
		{
			name:           "methods with spaces",
			method:         "PUT",
			allowedMethods: "GET, POST, PUT, DELETE",
			expected:       true,
		},
		{
			name:           "empty allowed methods",
			method:         "GET",
			allowedMethods: "",
			expected:       false,
		},
		{
			name:           "single allowed method",
			method:         "GET",
			allowedMethods: "GET",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middlewares.IsMethodAllowed(tt.method, tt.allowedMethods)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseAllowedOrigins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single origin",
			input:    "https://example.com",
			expected: []string{"https://example.com"},
		},
		{
			name:     "multiple origins",
			input:    "https://example.com, https://api.example.com",
			expected: []string{"https://example.com", "https://api.example.com"},
		},
		{
			name:     "origins with spaces",
			input:    "https://example.com, https://api.example.com, https://app.example.com",
			expected: []string{"https://example.com", "https://api.example.com", "https://app.example.com"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only spaces",
			input:    "   ,   ",
			expected: []string{},
		},
		{
			name:     "trailing comma",
			input:    "https://example.com,",
			expected: []string{"https://example.com"},
		},
		{
			name:     "wildcard origins",
			input:    "https://*.example.com, https://*.api.example.com",
			expected: []string{"https://*.example.com", "https://*.api.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middlewares.ParseAllowedOrigins(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		expected       bool
	}{
		{
			name:           "exact match",
			origin:         "https://example.com",
			allowedOrigins: []string{"https://example.com"},
			expected:       true,
		},
		{
			name:           "no match",
			origin:         "https://evil.com",
			allowedOrigins: []string{"https://example.com"},
			expected:       false,
		},
		{
			name:           "wildcard match",
			origin:         "https://api.example.com",
			allowedOrigins: []string{"https://*.example.com"},
			expected:       true,
		},
		{
			name:           "wildcard match - deep subdomain",
			origin:         "https://v1.api.example.com",
			allowedOrigins: []string{"https://*.example.com"},
			expected:       true,
		},
		{
			name:           "wildcard match without protocol",
			origin:         "https://api.example.com",
			allowedOrigins: []string{"*.example.com"},
			expected:       true,
		},
		{
			name:           "multiple origins with wildcard",
			origin:         "https://v1.api.example.com",
			allowedOrigins: []string{"https://app.example.com", "https://*.api.example.com"},
			expected:       true,
		},
		{
			name:           "multiple origins no match",
			origin:         "https://evil.com",
			allowedOrigins: []string{"https://app.example.com", "https://*.api.example.com"},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middlewares.IsOriginAllowed(tt.origin, tt.allowedOrigins)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchOrigin(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		pattern  string
		expected bool
	}{
		{
			name:     "exact match",
			origin:   "https://example.com",
			pattern:  "https://example.com",
			expected: true,
		},
		{
			name:     "exact match with different protocol",
			origin:   "http://example.com",
			pattern:  "https://example.com",
			expected: false,
		},
		{
			name:     "wildcard match - simple subdomain",
			origin:   "https://api.example.com",
			pattern:  "https://*.example.com",
			expected: true,
		},
		{
			name:     "wildcard match - deep subdomain",
			origin:   "https://v1.api.example.com",
			pattern:  "https://*.example.com",
			expected: true,
		},
		{
			name:     "wildcard match - http protocol",
			origin:   "http://api.example.com",
			pattern:  "http://*.example.com",
			expected: true,
		},
		{
			name:     "wildcard no match - protocol mismatch",
			origin:   "http://api.example.com",
			pattern:  "https://*.example.com",
			expected: false,
		},
		{
			name:     "wildcard no match - naked domain",
			origin:   "https://example.com",
			pattern:  "https://*.example.com",
			expected: false,
		},
		{
			name:     "wildcard no match - different domain",
			origin:   "https://api.other.com",
			pattern:  "https://*.example.com",
			expected: false,
		},
		{
			name:     "wildcard without protocol",
			origin:   "https://api.example.com",
			pattern:  "*.example.com",
			expected: true,
		},
		{
			name:     "wildcard without protocol - http",
			origin:   "http://api.example.com",
			pattern:  "*.example.com",
			expected: true,
		},
		{
			name:     "wildcard without protocol - no match",
			origin:   "https://api.other.com",
			pattern:  "*.example.com",
			expected: false,
		},
		{
			name:     "wildcard with multiple subdomains",
			origin:   "https://api.v1.example.com",
			pattern:  "https://*.example.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middlewares.MatchOrigin(tt.origin, tt.pattern)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCORSWithNilHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.CORS("https://example.com", "GET, POST", nil, nil))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "Content-Length, X-Request-ID",
		w.Header().Get("Access-Control-Expose-Headers"))
}

func TestCORSWithEmptyAllowedOrigins(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.CORS("", "GET, POST", nil, nil))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Request with origin
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Request without origin
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func BenchmarkCORS(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.CORS("https://*.example.com, https://app.example.com", "GET, POST", nil, nil))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://api.example.com")
		router.ServeHTTP(w, req)
	}
}

func stringPtr(s string) *string {
	return &s
}

func BenchmarkCORSExactMatch(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.CORS("https://example.com", "GET, POST", nil, nil))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		router.ServeHTTP(w, req)
	}
}

func BenchmarkCORSWildcardMatch(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.CORS("https://*.example.com", "GET, POST", nil, nil))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://api.example.com")
		router.ServeHTTP(w, req)
	}
}

func BenchmarkCORSNoOrigin(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.CORS("https://*.example.com", "GET, POST", nil, nil))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		// No Origin header set
		router.ServeHTTP(w, req)
	}
}
