package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func ApiRootRoute(rg *gin.RouterGroup, handler *handler.ApiRootHandler, rateLimiter *middleware.RateLimiter) {
	rg.GET("/", rateLimiter.Middleware(), handler.GetV1)
}
