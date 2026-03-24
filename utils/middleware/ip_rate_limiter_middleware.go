package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
	prefix      string        // prefix for redis keys
	rate        float64       // requests per second
	burst       int           // max burst
	window      time.Duration // fixed window duration
}

func NewRateLimiter(client *redis.Client, prefix string, rate float64, burst int) *RateLimiter {
	return &RateLimiter{
		redisClient: client,
		prefix:      prefix,
		rate:        rate,
		burst:       burst,
		window:      time.Minute, // Default to 1 minute window for simplicity
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If redis is not available, allow the request (fail-open)
		if rl.redisClient == nil {
			c.Next()
			return
		}

		userID, exists := c.Get("user_id")

		var identifier string
		if exists {
			identifier = "user:" + userID.(string)
		} else {
			identifier = "ip:" + c.ClientIP()
		}

		key := fmt.Sprintf("%s:%s", rl.prefix, identifier)
		ctx := c.Request.Context()

		// DEBUG: Log that we are checking Redis
		fmt.Printf("[DEBUG] Rate limiting check for key: %s\n", key)

		// Simple fixed-window rate limiting using Redis INCR
		count, err := rl.redisClient.Incr(ctx, key).Result()
		if err != nil {
			// On Redis error, we fail-open to not block the user
			c.Next()
			return
		}

		// Set expiration on the first request in the window
		if count == 1 {
			rl.redisClient.Expire(ctx, key, rl.window)
		}

		// Calculate the allowed limit for the window
		// Since rate is requests/sec, for 1 minute window it's rate * 60
		limit := int64(rl.rate * float64(rl.window.Seconds()))
		if limit == 0 {
			limit = int64(rl.burst)
		}

		if count > limit {
			retryAfter := int(rl.window.Seconds())
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
