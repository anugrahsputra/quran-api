package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	DefaultMaxBodySize = 1 << 20
)

func BodySizeLimit(maxSize int64) gin.HandlerFunc {
	if maxSize <= 0 {
		maxSize = DefaultMaxBodySize
	}

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}
