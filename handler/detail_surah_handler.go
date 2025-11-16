package handler

import (
	"net/http"
	"strconv"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/anugrahsputra/go-quran-api/utils/helper"
	"github.com/gin-gonic/gin"
)

type DetailSurahHandler struct {
	quranService service.IQuranService
}

func NewDetailSurahHandler(quranService service.IQuranService) *DetailSurahHandler {
	return &DetailSurahHandler{
		quranService: quranService,
	}
}

func (h *DetailSurahHandler) GetDetailSurah(c *gin.Context) {
	surahIDStr := c.Param("surah_id")
	surahID, err := strconv.Atoi(surahIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing surah_id",
		})
		return
	}

	// Validate surah_id range (1-114)
	if surahID < 1 || surahID > 114 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "surah_id must be between 1 and 114",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	logger.Infof(
		"HTTP %s %s | IP: %s | Params: surah_id=%d, page=%d, limit=%d | UA: %s",
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		surahID,
		page,
		limit,
		c.Request.UserAgent(),
	)

	data, totalVerses, totalPages, err := h.quranService.GetSurahDetail(c.Request.Context(), surahID, page, limit)
	if err != nil {
		logger.Errorf("Error fetching surah detail: %s", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: helper.SanitizeError(err),
		})
		return
	}

	logger.Infof("Request completed successfully - Path: %s", c.Request.URL.Path)

	c.JSON(http.StatusOK, dto.SurahDetailResp{
		Status:  http.StatusOK,
		Message: "success",
		Meta: dto.Meta{
			Total:      totalVerses,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
		Data: data,
	})
}
