package router

import (
	"github.com/anugrahsputra/quran-api/config"
	"github.com/anugrahsputra/quran-api/handler"
	"github.com/anugrahsputra/quran-api/repository"
	"github.com/anugrahsputra/quran-api/service"
	"github.com/gin-gonic/gin"
)

func SetupRoute() *gin.Engine {
	route := gin.Default()
	apiV1 := route.Group("/api/v1")
	cfg := config.LoadConfig()

	// List surah
	surahRepo := repository.NewQuranRepository(cfg)
	surahService := service.NewSurahService(surahRepo)
	surahHandler := handler.NewSurahHandler(surahService)

	SurahRoute(apiV1, surahHandler)

	return route
}
