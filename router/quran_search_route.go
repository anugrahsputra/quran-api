package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func NewQuranSearchRoute(g *gin.RouterGroup, h *handler.QuranSearchHandler, rl *middleware.RateLimiter) {
	g.GET("/search", rl.Middleware(), h.Search)
}
