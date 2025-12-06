package router

import (
	"log"
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

func SetupRoute() *gin.Engine {
	// Set Gin mode based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		// If GIN_MODE is not set, check ENV variable
		env := os.Getenv("ENV")
		if env == "production" || env == "prod" {
			gin.SetMode(gin.ReleaseMode)
		} else {
			// Default to debug mode for development
			gin.SetMode(gin.DebugMode)
		}
	} else {
		// Use GIN_MODE directly if set
		gin.SetMode(ginMode)
	}

	route := gin.Default()
	cfg := config.LoadConfig()

	// Add global middlewares (order matters!)
	route.Use(middleware.Recovery())                // Panic recovery first
	route.Use(middleware.RequestID())               // Request ID for tracing
	route.Use(middleware.SecurityHeaders())         // Security headers
	route.Use(middleware.Timeout(30 * time.Second)) // 30 second timeout for all requests
	route.Use(middleware.BodySizeLimit(1 << 20))    // 1MB body size limit

	// Health check routes (no rate limiting, accessible at root level)
	searchRepo, err := repository.NewQuranSearchRepository(cfg.SearchIndexPath)
	if err != nil {
		log.Fatalf("failed to create search repository: %v", err)
	}
	healthHandler := handler.NewHealthHandler(searchRepo)
	HealthRoute(route.Group(""), healthHandler)

	apiV1 := route.Group("/api/v1")
	rateLimiter := middleware.NewRateLimiter(rate.Every(time.Minute/5), 5)

	quranRepo := repository.NewQuranRepository(cfg)
	quranService := service.NewQuranService(quranRepo)
	// surah list
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

	searchService := service.NewQuranSearchService(quranRepo, searchRepo)
	searchHandler := handler.NewQuranSearchHandler(searchService)
	NewQuranSearchRoute(apiV1, searchHandler, rateLimiter)

	// Admin routes (for administrative operations like reindexing)
	adminHandler := handler.NewAdminHandler(searchService)
	AdminRoute(apiV1, adminHandler, rateLimiter)

	return route
}
