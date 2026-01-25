package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters        map[string]*rate.Limiter
	lastAccess      map[string]time.Time
	mu              *sync.Mutex
	r               rate.Limit
	burst           int
	cleanupInterval time.Duration
	lastCleanup     time.Time
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters:        make(map[string]*rate.Limiter),
		lastAccess:      make(map[string]time.Time),
		mu:              &sync.Mutex{},
		r:               r,
		burst:           b,
		cleanupInterval: 5 * time.Minute,
		lastCleanup:     time.Now(),
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastCleanup) > rl.cleanupInterval {
		rl.cleanup(now)
		rl.lastCleanup = now
	}

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.burst)
		rl.limiters[key] = limiter
	}

	rl.lastAccess[key] = now
	return limiter
}

func (rl *RateLimiter) cleanup(now time.Time) {
	cutoff := now.Add(-rl.cleanupInterval * 2)

	for key, lastAccess := range rl.lastAccess {
		if lastAccess.Before(cutoff) {
			delete(rl.limiters, key)
			delete(rl.lastAccess, key)
		}
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")

		var key string
		if exists {
			key = "user:" + userID.(string)
		} else {
			key = "ip:" + c.ClientIP()
		}

		limiter := rl.getLimiter(key)
		if !limiter.Allow() {
			retryAfter := 1
			if rl.r > 0 {
				retryAfter = int(1 / float64(rl.r))
				if retryAfter < 1 {
					retryAfter = 1
				}
			}

			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Status:  http.StatusTooManyRequests,
				Message: "Too Many Requests",
			})
			return
		}

		c.Next()
	}
}
