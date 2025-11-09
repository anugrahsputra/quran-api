package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout middleware sets a timeout for request processing
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			// Request completed within timeout
		case <-ctx.Done():
			// Request timed out
			c.AbortWithStatus(408) // Request Timeout
		}
	}
}
