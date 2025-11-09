package middleware

import (
	"net"
	"net/http"
	"sync"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       *sync.Mutex
	r        rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		mu:       &sync.Mutex{},
		r:        r,
		burst:    b,
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.burst)
		rl.limiters[key] = limiter
	}
	return limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")

		var key string
		if exists {
			key = "user:" + userID.(string)
		} else {
			key = "ip:" + getClientIP(c.Request)
		}

		limiter := rl.getLimiter(key)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Status:  http.StatusTooManyRequests,
				Message: "Too Many Request;",
			})
			return
		}

		c.Next()
	}
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}
