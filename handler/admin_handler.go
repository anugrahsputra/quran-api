package handler

import (
	"log"
	"net/http"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	searchAyahService service.SearchAyahService
}

func NewAdminHandler(searchAyahService service.SearchAyahService) *AdminHandler {
	return &AdminHandler{
		searchAyahService: searchAyahService,
	}
}

func (h *AdminHandler) Reindex(c *gin.Context) {
	go func() {
		log.Println("Starting reindexing process via API...")
		if err := h.searchAyahService.IndexQuran(); err != nil {
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
