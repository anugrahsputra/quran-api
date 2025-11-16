package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func SurahRoute(r *gin.RouterGroup, h *handler.SurahHandler, rl *middleware.RateLimiter) {
	surahGroup := r.Group("/surah", rl.Middleware())
	{
		surahGroup.GET("/", h.GetListSurah)
	}
}
