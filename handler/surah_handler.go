package handler

import (
	"net/http"
	"strconv"

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
	logger.Infof(
		"HTTP request received - Method: %s, Path: %s, RemoteAddr: %s, UserAgent: %s",
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		c.Request.UserAgent(),
	)

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

	c.JSON(http.StatusOK, dto.SurahListResp{
		Status:  http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

func (s *SurahHandler) GetDetailSurah(c *gin.Context) {
	logger.Infof(
		"HTTP %s %s | IP: %s | Params: surah_id=%s, start=%s, limit=%s | UA: %s",
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		c.Query("surah_id"),
		c.Query("start"),
		c.Query("limit"),
		c.Request.UserAgent(),
	)

	startStr := c.DefaultQuery("start", "0")
	limitStr := c.DefaultQuery("limit", "10")
	surahIDStr := c.Query("surah_id")

	surahID, err := strconv.Atoi(surahIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing surah_id",
		})
		return
	}

	start, _ := strconv.Atoi(startStr)
	pageLimit, _ := strconv.Atoi(limitStr)

	response, err := s.surahService.GetSurahDetail(c.Request.Context(), surahID, start, pageLimit)
	if err != nil {
		logger.Errorf("Error fetching surah detail: %s", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	logger.Infof("Request completed successfully - Path: %s", c.Request.URL.Path)

	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    response,
	})
}
