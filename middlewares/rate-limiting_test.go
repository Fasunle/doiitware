package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fasunle/doiitware/helpers"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestRateLimitMiddlewareAllowsWhitelistedIP(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitMiddleware(helpers.NewIPRateLimiter(0, 0), []string{"127.0.0.1/32"}))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.RemoteAddr = "127.0.0.1:1234"

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected whitelisted request to pass through, got status %d", recorder.Code)
	}
}

func TestRateLimitMiddlewareBlocksExcessRequests(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitMiddleware(helpers.NewIPRateLimiter(rate.Limit(0), 0), nil))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.RemoteAddr = "203.0.113.10:1234"

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("expected rate-limited request to return 429, got %d", recorder.Code)
	}

	if body := recorder.Body.String(); body == "" {
		t.Fatalf("expected rate-limited response body")
	}
}
