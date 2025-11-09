package router

import (
	"log"
	"time"

	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/repository"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
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
	DetailSurahRoute(apiV1, surahHandler, rateLimiter)

	prayerTimeRepo := repository.NewPrayerTimeRepository(cfg)
	prayerTimeService := service.NewPrayerTimeService(prayerTimeRepo)
	prayerTimeHandler := handler.NewPrayerTimeHandler(prayerTimeService)
	PrayerTimeRoute(apiV1, prayerTimeHandler, rateLimiter)

	searchRepo, err := repository.NewSearchRepository()
	if err != nil {
		log.Fatalf("failed to create search repository: %v", err)
	}
	searchService := service.NewSearchService(surahRepo, searchRepo)
	searchHandler := handler.NewSearchHandler(searchService)
	NewSearchRoute(apiV1, searchHandler, rateLimiter)

	return route
}
