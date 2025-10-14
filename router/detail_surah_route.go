package router

import (
	"github.com/anugrahsputra/quran-api/handler"
	"github.com/anugrahsputra/quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func DetailSurahRoute(r *gin.RouterGroup, surahHandler *handler.SurahHandler, rateLimiter *middleware.RateLimiter) {
	surahGroup := r.Group("/surah/:surah_id", rateLimiter.Middleware())
	{
		surahGroup.GET("/", surahHandler.GetDetailSurah)
	}
}
