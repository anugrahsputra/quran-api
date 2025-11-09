package router

import (
	"github.com/anugrahsputra/go-quran-api/handler"
		"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
)

func PrayerTimeRoute(r *gin.RouterGroup, prayerTimeHandler *handler.PrayerTimeHandler, rateLimiter *middleware.RateLimiter) {
	prayerTimeGroup := r.Group("/prayer-time", rateLimiter.Middleware())
	{
		prayerTimeGroup.GET("/", prayerTimeHandler.GetPrayerTime)
	}
}
