package handler

import (
	"net/http"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/anugrahsputra/go-quran-api/utils/helper"
	"github.com/gin-gonic/gin"
)

type PrayerTimeHandler struct {
	prayerTimeService service.IPrayerTimeService
}

func NewPrayerTimeHandler(prayerTimeService service.IPrayerTimeService) *PrayerTimeHandler {
	return &PrayerTimeHandler{
		prayerTimeService: prayerTimeService,
	}
}

func (s *PrayerTimeHandler) GetPrayerTime(c *gin.Context) {
	// Logging removed for now - can be added back with shared logger if needed

	city := c.DefaultQuery("address", "Jakarta")
	timezone := c.DefaultQuery("timezonestring", "Asia/Jakarta")

	response, err := s.prayerTimeService.GetPrayerTime(c.Request.Context(), city, timezone)
	if err != nil {
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
