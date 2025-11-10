package middleware

import (
	"os"
	"strings"

	"github.com/anugrahsputra/go-quran-api/utils/helper"
	"github.com/gin-gonic/gin"
)

// CORS middleware for handling cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		isProduction := helper.IsProduction()

		// In production, use whitelist from environment variable
		// In development, allow all origins for easier testing
		if isProduction {
			allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
			if allowedOrigins == "" {
				// Default to empty in production if not set (most restrictive)
				c.Header("Access-Control-Allow-Origin", "")
			} else {
				// Check if origin is in whitelist
				origins := strings.Split(allowedOrigins, ",")
				allowed := false
				for _, allowedOrigin := range origins {
					if strings.TrimSpace(allowedOrigin) == origin {
						allowed = true
						break
					}
				}
				if allowed {
					c.Header("Access-Control-Allow-Origin", origin)
				} else {
					// Origin not in whitelist, deny
					c.Header("Access-Control-Allow-Origin", "")
				}
			}
		} else {
			// Development: Allow all origins
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
