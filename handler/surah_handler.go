package handler

import (
	"net/http"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/anugrahsputra/go-quran-api/utils/helper"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("handler")

type SurahHandler struct {
	quranService service.IQuranService
}

func NewSurahHandler(surahService service.IQuranService) *SurahHandler {
	return &SurahHandler{
		quranService: surahService,
	}
}

func (h *SurahHandler) GetListSurah(c *gin.Context) {
	logger.Infof(
		"HTTP request received - Method: %s, Path: %s, RemoteAddr: %s, UserAgent: %s",
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	response, err := h.quranService.GetListSurah(c.Request.Context())
	if err != nil {
		logger.Errorf("HTTP request failed - Method: %s, Path: %s, Error: %s",
			c.Request.Method, c.Request.URL.Path, err.Error())
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: helper.SanitizeError(err),
		})
		return
	}

	logger.Infof("HTTP request completed successfully - Method: %s, Path: %s, Status: %d",
		c.Request.Method, c.Request.URL.Path, http.StatusOK)

	c.JSON(http.StatusOK, dto.SurahListResp{
		Status:  http.StatusOK,
		Message: "success",
		Data:    response,
	})
}
