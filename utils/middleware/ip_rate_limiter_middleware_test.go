package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should allow request within limit", func(t *testing.T) {
		rl := NewRateLimiter(rate.Limit(10), 1)
		router := gin.New()
		router.Use(rl.Middleware())
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should block request exceeding limit and include Retry-After header", func(t *testing.T) {
		rl := NewRateLimiter(rate.Every(10*time.Second), 1)
		router := gin.New()
		router.Use(rl.Middleware())
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
		assert.Equal(t, "10", w2.Header().Get("Retry-After"))
		assert.Contains(t, w2.Body.String(), "Too Many Requests")
	})

	t.Run("should use 1 as minimum Retry-After if rate is high", func(t *testing.T) {
		rl := NewRateLimiter(rate.Limit(100), 1)
		router := gin.New()
		router.Use(rl.Middleware())
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
		assert.Equal(t, "1", w2.Header().Get("Retry-After"))
	})
}
