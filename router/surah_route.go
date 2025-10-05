package router

import (
	"github.com/anugrahsputra/quran-api/handler"
	"github.com/gin-gonic/gin"
)

func SurahRoute(api *gin.RouterGroup, surahHandler *handler.SurahHandler) {
	api.GET("/surahs", surahHandler.GetListSurah)
}
