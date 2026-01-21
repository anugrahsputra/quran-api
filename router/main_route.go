package router

import (
	"os"
	"time"

	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/handler"
	"github.com/anugrahsputra/go-quran-api/repository"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRoute(
	cfg *config.Config,
	quranRepo repository.IQuranRepository,
	searchRepo repository.QuranSearchRepository,
	searchService service.QuranSearchService,
) *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		env := os.Getenv("ENV")
		if env == "production" || env == "prod" {
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}
	} else {
		gin.SetMode(ginMode)
	}

	route := gin.Default()

	route.Use(middleware.Recovery())
	route.Use(middleware.RequestID())
	route.Use(middleware.SecurityHeaders())
	route.Use(middleware.Timeout(30 * time.Second))
	route.Use(middleware.BodySizeLimit(1 << 20))

	healthHandler := handler.NewHealthHandler(searchRepo)
	HealthRoute(route.Group(""), healthHandler)

	apiV1 := route.Group("/api/v1")
	rateLimiter := middleware.NewRateLimiter(rate.Limit(1), 10)

	quranService := service.NewQuranService(quranRepo)
	surahHandler := handler.NewSurahHandler(quranService)
	detailSurahHandler := handler.NewDetailSurahHandler(quranService)
	detailAyahHandler := handler.NewDetailAyahHandler(quranService)
	SurahRoute(apiV1, surahHandler, rateLimiter)
	DetailSurahRoute(apiV1, detailSurahHandler, rateLimiter)
	DetailAyahRoute(apiV1, detailAyahHandler, rateLimiter)

	prayerTimeRepo := repository.NewPrayerTimeRepository(cfg)
	prayerTimeService := service.NewPrayerTimeService(prayerTimeRepo)
	prayerTimeHandler := handler.NewPrayerTimeHandler(prayerTimeService)
	PrayerTimeRoute(apiV1, prayerTimeHandler, rateLimiter)

	searchHandler := handler.NewQuranSearchHandler(searchService)
	NewQuranSearchRoute(apiV1, searchHandler, rateLimiter)

	adminHandler := handler.NewAdminHandler(searchService)
	AdminRoute(apiV1, adminHandler, rateLimiter)

	return route
}