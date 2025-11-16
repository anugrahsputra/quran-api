package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func PrayerTimeRoute(r *gin.RouterGroup, h *handler.PrayerTimeHandler, rl *middleware.RateLimiter) {
	prayerTimeGroup := r.Group("/prayer-time", rl.Middleware())
	{
		prayerTimeGroup.GET("/", h.GetPrayerTime)
	}
}
