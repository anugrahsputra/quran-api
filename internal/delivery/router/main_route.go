package router

import (
	"os"
	"time"

	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/internal/delivery/handler"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
	"github.com/anugrahsputra/go-quran-api/internal/repository"
	"github.com/anugrahsputra/go-quran-api/internal/service"
	"github.com/anugrahsputra/go-quran-api/utils/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func wireHealth(searchRepo domain.QuranSearchRepository) *handler.HealthHandler {
	return handler.NewHealthHandler(searchRepo)
}

func wireApiRootRoute(rc *redis.Client) *handler.ApiRootHandler {
	apiRootRepo := repository.NewApiRootRepository()
	apiRootService := service.NewApiRootService(apiRootRepo)
	return handler.NewApiRootHandler(apiRootService)
}

func wireSurahRoutes(surahRepo domain.SurahRepository, ayahRepo domain.AyahRepository, rc *redis.Client) (*handler.SurahHandler, *handler.DetailSurahHandler, *handler.DetailAyahHandler) {
	surahService := service.NewSurahService(surahRepo, rc)
	ayahService := service.NewAyahService(ayahRepo, rc)
	return handler.NewSurahHandler(surahService), handler.NewDetailSurahHandler(surahService), handler.NewDetailAyahHandler(ayahService)
}

func wirePrayerTime(cfg *config.Config) *handler.PrayerTimeHandler {
	prayerTimeRepo := repository.NewPrayerTimeRepository(cfg)
	prayerTimeService := service.NewPrayerTimeService(prayerTimeRepo)
	return handler.NewPrayerTimeHandler(prayerTimeService)
}

func wireQuranSearch(searchService service.IQuranSearchService) (*handler.QuranSearchHandler, *handler.AdminHandler) {
	return handler.NewQuranSearchHandler(searchService), handler.NewAdminHandler(searchService)
}

type RouterDeps struct {
	Cfg           *config.Config
	SurahRepo     domain.SurahRepository
	AyahRepo      domain.AyahRepository
	SearchRepo    domain.QuranSearchRepository
	SearchService service.IQuranSearchService
	RedisClient   *redis.Client
}

func SetupRoute(deps RouterDeps) *gin.Engine {
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

	healthHandler := wireHealth(deps.SearchRepo)
	HealthRoute(route.Group(""), healthHandler)

	api := route.Group("/api")
	rateLimiter := middleware.NewRateLimiter(deps.RedisClient, "ratelimit:quran-api", 2.0, 120)

	apiRootHandler := wireApiRootRoute(deps.RedisClient)
	ApiRootRoute(api, apiRootHandler, rateLimiter)

	apiV1 := api.Group("/v1")

	surahHandler, detailSurahHandler, detailAyahHandler := wireSurahRoutes(deps.SurahRepo, deps.AyahRepo, deps.RedisClient)
	SurahRoute(apiV1, surahHandler, rateLimiter)
	DetailSurahRoute(apiV1, detailSurahHandler, rateLimiter)
	DetailAyahRoute(apiV1, detailAyahHandler, rateLimiter)

	prayerTimeHandler := wirePrayerTime(deps.Cfg)
	PrayerTimeRoute(apiV1, prayerTimeHandler, rateLimiter)

	searchHandler, adminHandler := wireQuranSearch(deps.SearchService)
	NewQuranSearchRoute(apiV1, searchHandler, rateLimiter)
	AdminRoute(apiV1, adminHandler, rateLimiter)

	return route
}
