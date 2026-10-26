package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimit(t *testing.T) {
	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		t.Fatal(err)
	}

	calls := 0
	handler := func(c *gin.Context) {
		calls++
		c.Status(http.StatusNoContent)
	}

	limiter := RateLimit()
	router.POST("/auth/login", limiter, handler)
	router.POST("/auth/register", limiter, handler)

	request := func(ip, path, forwardedIP string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.RemoteAddr = ip + ":12345"
		req.Header.Set("X-Forwarded-For", forwardedIP)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		return recorder
	}

	for i := 0; i < 5; i++ {
		path := "/auth/login"
		if i%2 == 0 {
			path = "/auth/register"
		}
		response := request("192.0.2.1", path, "203.0.113.1")
		if response.Code != http.StatusNoContent {
			t.Fatalf("request %d: got %d", i+1, response.Code)
		}
	}

	response := request("192.0.2.1", "/auth/login", "203.0.113.2")
	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("over limit: got %d", response.Code)
	}
	if calls != 5 {
		t.Fatalf("blocked request reached handler: calls=%d", calls)
	}
	retryAfter, err := strconv.Atoi(response.Header().Get("Retry-After"))
	if err != nil || retryAfter < 1 || retryAfter > 60 {
		t.Fatalf("invalid Retry-After: %q", response.Header().Get("Retry-After"))
	}

	response = request("192.0.2.2", "/auth/register", "")
	if response.Code != http.StatusNoContent {
		t.Fatalf("different IP: got %d", response.Code)
	}
}
