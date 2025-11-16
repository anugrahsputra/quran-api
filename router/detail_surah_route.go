package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func DetailSurahRoute(r *gin.RouterGroup, h *handler.DetailSurahHandler, rl *middleware.RateLimiter) {
	surahGroup := r.Group("/surah/:surah_id", rl.Middleware())
	{
		surahGroup.GET("/", h.GetDetailSurah)
	}
}
