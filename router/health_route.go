package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/gin-gonic/gin"
)

func HealthRoute(g *gin.RouterGroup, h *handler.HealthHandler) {
	g.GET("/health", h.HealthCheck)
	g.GET("/health/live", h.LivenessCheck)
	g.GET("/health/ready", h.ReadinessCheck)
}
