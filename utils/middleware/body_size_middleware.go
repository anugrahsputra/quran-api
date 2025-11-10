package middleware

import (
	"net/http"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
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
		// Only check body size for methods that typically have bodies
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.ContentLength > maxSize {
				c.JSON(http.StatusRequestEntityTooLarge, dto.ErrorResponse{
					Status:  http.StatusRequestEntityTooLarge,
					Message: "Request body too large",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
