package handler

import (
	"github.com/anugrahsputra/quran-api/domain/dto"
	"github.com/anugrahsputra/quran-api/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
	logger.Infof(
		"HTTP request received - Method: %s, Path: %s, RemoteAddr: %s, UserAgent: %s",
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		c.Query("address"),
		c.Query("timezonestring"),
		c.Request.UserAgent(),
	)

	city := c.DefaultQuery("address", "Jakarta")
	timezone := c.DefaultQuery("timezonestring", "Asia/Jakarta")

	response, err := s.prayerTimeService.GetPrayerTime(c.Request.Context(), city, timezone)
	if err != nil {
		logger.Errorf("HTTP request failed - Method: %s, Path: %s, Error: %s", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	logger.Infof("HTTP request completed successfully - Method: %s, Path: %s, Status: %d",
		c.Request.Method, c.Request.URL.Path, http.StatusOK)

	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    response,
	})
}
