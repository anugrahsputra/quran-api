package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func DetailAyahRoute(r *gin.RouterGroup, h *handler.DetailAyahHandler, rl *middleware.RateLimiter) {
	ayahGroup := r.Group("/ayah/:ayah_id", rl.Middleware())
	{
		ayahGroup.GET("/", h.GetDetailAyah)
	}
}
