package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func NewSearchAyahRoute(g *gin.RouterGroup, h *handler.SearchAyahHandler, rl *middleware.RateLimiter) {
	g.GET("/search", rl.Middleware(), h.Search)
}
