package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/gin-gonic/gin"
)

func HealthRoute(g *gin.RouterGroup, healthHandler *handler.HealthHandler) {
	// Health check endpoint (can include dependency checks)
	g.GET("/health", healthHandler.HealthCheck)

	// Liveness probe (simple - just checks if service is running)
	g.GET("/health/live", healthHandler.LivenessCheck)

	// Readiness probe (checks if service is ready to accept traffic)
	g.GET("/health/ready", healthHandler.ReadinessCheck)
}
