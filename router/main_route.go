package router

import (
	"time"

	"github.com/anugrahsputra/quran-api/config"
	"github.com/anugrahsputra/quran-api/handler"
	"github.com/anugrahsputra/quran-api/repository"
	"github.com/anugrahsputra/quran-api/service"
	"github.com/anugrahsputra/quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRoute() *gin.Engine {
	route := gin.Default()
	cfg := config.LoadConfig()
	apiV1 := route.Group("/api/v1")
	rateLimiter := middleware.NewRateLimiter(rate.Every(time.Minute/5), 5)

	surahRepo := repository.NewQuranRepository(cfg)
	surahService := service.NewSurahService(surahRepo)
	surahHandler := handler.NewSurahHandler(surahService)
	SurahRoute(apiV1, surahHandler, rateLimiter)

	return route
}
