package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func SurahRoute(r *gin.RouterGroup, surahHandler *handler.SurahHandler, rateLimiter *middleware.RateLimiter) {
	surahGroup := r.Group("/surah", rateLimiter.Middleware())
	{
		surahGroup.GET("/", surahHandler.GetListSurah)
		surahGroup.GET("/detail", surahHandler.GetDetailSurah)
	}
}
