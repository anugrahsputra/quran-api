package handler

import (
	"net/http"
	"strconv"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/anugrahsputra/go-quran-api/utils/helper"
	"github.com/gin-gonic/gin"
)

type SearchAyahHandler struct {
	searchAyahService service.SearchAyahService
}

func NewSearchAyahHandler(searchAyahService service.SearchAyahService) *SearchAyahHandler {
	return &SearchAyahHandler{searchAyahService: searchAyahService}
}

func (h *SearchAyahHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.SearchResponse{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Query parameter 'q' is required",
			Meta: dto.Meta{
				Total:      0,
				Page:       1,
				Limit:      10,
				TotalPages: 0,
			},
		})
		return
	}

	// Parse pagination parameters with defaults
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	ayahs, total, err := h.searchAyahService.Search(query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.SearchResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: helper.SanitizeError(err),
			Meta: dto.Meta{
				Total:      0,
				Page:       page,
				Limit:      limit,
				TotalPages: 0,
			},
		})
		return
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit // Ceiling division
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, dto.SearchResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Success",
		Meta: dto.Meta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
		Data: ayahs,
	})
}
