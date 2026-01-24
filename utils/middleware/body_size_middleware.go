package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultMaxBodySize is 1MB
	DefaultMaxBodySize = 1 << 20 // 1MB
)

// BodySizeLimit middleware limits the size of request body
func BodySizeLimit(maxSize int64) gin.HandlerFunc {
	if maxSize <= 0 {
		maxSize = DefaultMaxBodySize
	}

	return func(c *gin.Context) {
		// Set a limit on the request body size
		// This handles chunked encoding and prevents reading more than maxSize
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}
