package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCustomRecoveryMiddlewareReturnsSafeErrorResponse(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CustomRecoveryMiddleware())
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "req-123")
		c.Next()
	})
	router.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	request.RemoteAddr = "127.0.0.1:1234"

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected recovery middleware to return 500, got %d", recorder.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON error response: %v", err)
	}

	if body["error"] != "internal_server_error" {
		t.Fatalf("expected internal_server_error response, got %v", body["error"])
	}

	if body["request_id"] != "req-123" {
		t.Fatalf("expected request_id to be propagated, got %v", body["request_id"])
	}
}
