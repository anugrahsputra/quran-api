package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/gin-gonic/gin"
)

func HealthRoute(g *gin.RouterGroup, h *handler.HealthHandler) {
	// Health check endpoint (can include dependency checks)
	g.GET("/health", h.HealthCheck)

	// Liveness probe (simple - just checks if service is running)
	g.GET("/health/live", h.LivenessCheck)

	// Readiness probe (checks if service is ready to accept traffic)
	g.GET("/health/ready", h.ReadinessCheck)
}
