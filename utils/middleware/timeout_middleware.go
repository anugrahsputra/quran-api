package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout middleware sets a timeout for request processing by wrapping the request context
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the request context with the timed-out one
		c.Request = c.Request.WithContext(ctx)

		// Continue processing
		c.Next()

		// If the context timed out during execution, abort with 408
		if ctx.Err() == context.DeadlineExceeded {
			c.AbortWithStatus(408)
		}
	}
}
