package handler

import (
	"net/http"
	"time"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/repository"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	searchRepo repository.QuranSearchRepository
}

func NewHealthHandler(searchRepo repository.QuranSearchRepository) *HealthHandler {
	return &HealthHandler{
		searchRepo: searchRepo,
	}
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Checks    map[string]HealthCheck `json:"checks"`
}

type HealthCheck struct {
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	ResponseTime string `json:"response_time,omitempty"`
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	checks := make(map[string]HealthCheck)
	overallStatus := "healthy"

	// Check search index
	indexCheck := h.checkSearchIndex()
	checks["search_index"] = indexCheck
	if indexCheck.Status != "healthy" {
		overallStatus = "degraded"
	}

	// Check external APIs (optional - can be slow, so we'll do quick checks)
	// We'll skip this for now to keep health check fast, but can add if needed

	statusCode := http.StatusOK
	if overallStatus == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	})
}

func (h *HealthHandler) checkSearchIndex() HealthCheck {
	start := time.Now()

	// Check if index is healthy
	if !h.searchRepo.IsHealthy() {
		responseTime := time.Since(start)
		return HealthCheck{
			Status:       "unhealthy",
			Message:      "Search index is not accessible",
			ResponseTime: responseTime.String(),
		}
	}

	// Get document count to verify index has data
	docCount, err := h.searchRepo.GetDocCount()
	responseTime := time.Since(start)

	if err != nil {
		return HealthCheck{
			Status:       "unhealthy",
			Message:      "Failed to get document count: " + err.Error(),
			ResponseTime: responseTime.String(),
		}
	}

	message := "Search index is accessible"
	if docCount == 0 {
		message += " (no documents indexed)"
	} else {
		message += " (indexed documents available)"
	}

	return HealthCheck{
		Status:       "healthy",
		Message:      message,
		ResponseTime: responseTime.String(),
	}
}

// Readiness check - more comprehensive than health check
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	checks := make(map[string]HealthCheck)
	overallStatus := "ready"

	// Check search index
	indexCheck := h.checkSearchIndex()
	checks["search_index"] = indexCheck
	if indexCheck.Status != "healthy" {
		overallStatus = "not_ready"
	}

	statusCode := http.StatusOK
	if overallStatus == "not_ready" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, dto.Response{
		Status:  statusCode,
		Message: overallStatus,
		Data: map[string]interface{}{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"checks":    checks,
		},
	})
}

// Liveness check - simple check to see if service is alive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "alive",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make(map[string]HealthCheck),
	})
}
