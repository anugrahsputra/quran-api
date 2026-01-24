package middleware

import (
	"net/http"
	"os"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/gin-gonic/gin"
)

// AdminAuth middleware checks for a valid API key in the Authorization header
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := os.Getenv("ADMIN_API_KEY")

		// If no key is configured, default to a secure posture (deny all)
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Admin API key not configured on server",
			})
			return
		}

		clientKey := c.GetHeader("X-Admin-Token")
		if clientKey == "" || clientKey != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Invalid or missing admin token",
			})
			return
		}

		c.Next()
	}
}
