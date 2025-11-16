package handler

import (
	"log"
	"net/http"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	searchService service.SearchService
}

func NewAdminHandler(searchService service.SearchService) *AdminHandler {
	return &AdminHandler{
		searchService: searchService,
	}
}

func (h *AdminHandler) Reindex(c *gin.Context) {
	go func() {
		log.Println("Starting reindexing process via API...")
		if err := h.searchService.IndexQuran(); err != nil {
			log.Printf("Reindexing failed: %v", err)
		} else {
			log.Println("Reindexing completed successfully")
		}
	}()

	c.JSON(http.StatusAccepted, dto.Response{
		Status:  http.StatusAccepted,
		Message: "Reindexing started in background",
		Data: map[string]interface{}{
			"message": "The reindexing process has been started. This may take several minutes to complete.",
		},
	})
}
