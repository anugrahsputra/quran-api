package handler

import (
	"net/http"

	"github.com/anugrahsputra/quran-api/domain/dto"
	"github.com/anugrahsputra/quran-api/service"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("handler")

type SurahHandler struct {
	surahService service.ISurahService
}

func NewSurahHandler(surahService service.ISurahService) *SurahHandler {
	return &SurahHandler{
		surahService: surahService,
	}
}

func (s *SurahHandler) GetListSurah(c *gin.Context) {
	logger.Infof("HTTP request received - Method: %s, Path: %s, RemoteAddr: %s, UserAgent: %s",
		c.Request.Method, c.Request.URL.Path, c.ClientIP(), c.Request.UserAgent())

	response, err := s.surahService.GetListSurah(c.Request.Context())
	if err != nil {
		logger.Errorf("HTTP request failed - Method: %s, Path: %s, Error: %s",
			c.Request.Method, c.Request.URL.Path, err.Error())
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	logger.Infof("HTTP request completed successfully - Method: %s, Path: %s, Status: %d",
		c.Request.Method, c.Request.URL.Path, http.StatusOK)

	c.JSON(http.StatusOK, response)
}
