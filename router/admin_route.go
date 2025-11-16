package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoute(g *gin.RouterGroup, h *handler.AdminHandler, rl *middleware.RateLimiter) {
	g.POST("/reindex", rl.Middleware(), h.Reindex)
}
