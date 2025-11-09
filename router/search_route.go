package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func NewSearchRoute(g *gin.RouterGroup, searchHandler *handler.SearchHandler, rateLimiter *middleware.RateLimiter) {
	g.GET("/search", rateLimiter.Middleware(), searchHandler.Search)
}
