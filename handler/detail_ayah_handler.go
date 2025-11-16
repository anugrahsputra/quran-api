package handler

import (
	"net/http"
	"strconv"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/anugrahsputra/go-quran-api/utils/helper"
	"github.com/gin-gonic/gin"
)

type DetailAyahHandler struct {
	detailAyahService service.IQuranService
}

func NewDetailAyahHandler(das service.IQuranService) *DetailAyahHandler {
	return &DetailAyahHandler{
		detailAyahService: das,
	}
}

func (s *DetailAyahHandler) GetDetailAyah(c *gin.Context) {
	ayahIdStr := c.Param("ayah_id")
	ayahID, err := strconv.Atoi(ayahIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing ayah_id",
		})
		return
	}

	logger.Infof(
		"HTTP %s %s | IP: %s | Params: ayah_id=%d | UA: %s",
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		ayahID,
		c.Request.UserAgent(),
	)

	response, err := s.detailAyahService.GetDetailAyah(c.Request.Context(), ayahID)
	if err != nil {
		logger.Errorf("Error fetching ayah detail: %s", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: helper.SanitizeError(err),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    response,
	})

}
