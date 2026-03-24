package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to run miniredis: %v", err)
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	t.Run("should allow request within limit", func(t *testing.T) {
		s.FlushAll()
		// 10 requests per second = 600 requests per minute
		rl := NewRateLimiter(rdb, 10, 1)
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
		s.FlushAll()
		// 1 request per 60 seconds (rate=1/60)
		rl := NewRateLimiter(rdb, 1.0/60.0, 1)
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
		assert.Equal(t, "60", w2.Header().Get("Retry-After"))
		assert.Contains(t, w2.Body.String(), "Too Many Requests")
	})

	t.Run("should fail-open if redis is nil", func(t *testing.T) {
		rl := NewRateLimiter(nil, 10, 1)
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
}
