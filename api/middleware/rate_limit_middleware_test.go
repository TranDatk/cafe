package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Allowed Requests", func(t *testing.T) {
		r := gin.New()
		// 100 requests per second, burst 100
		r.Use(RateLimitMiddleware(100, 100))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:1234"
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("Blocked Requests", func(t *testing.T) {
		r := gin.New()
		// 1 request per second, burst 1
		r.Use(RateLimitMiddleware(1, 1))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// First request should pass
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		req1.RemoteAddr = "192.168.1.2:1234"
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second request should be blocked (429)
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		req2.RemoteAddr = "192.168.1.2:1234"
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	})

	t.Run("IP Isolation", func(t *testing.T) {
		r := gin.New()
		// 1 request per second, burst 1
		r.Use(RateLimitMiddleware(1, 1))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// IP 1: First request pass
		req1, _ := http.NewRequest("GET", "/test", nil)
		req1.RemoteAddr = "1.1.1.1:1234"
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// IP 1: Second request block
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req1)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)

		// IP 2: First request should still pass
		req2, _ := http.NewRequest("GET", "/test", nil)
		req2.RemoteAddr = "2.2.2.2:1234"
		w3 := httptest.NewRecorder()
		r.ServeHTTP(w3, req2)
		assert.Equal(t, http.StatusOK, w3.Code)
	})
}
