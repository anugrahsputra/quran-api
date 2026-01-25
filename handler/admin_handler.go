package handler

import (
	"log"
	"net/http"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	searchAyahService service.IQuranSearchService
}

func NewAdminHandler(sas service.IQuranSearchService) *AdminHandler {
	return &AdminHandler{
		searchAyahService: sas,
	}
}

func (h *AdminHandler) Reindex(c *gin.Context) {
	go func() {
		log.Println("Starting reindexing process via API...")
		if err := h.searchAyahService.IndexQuran(); err != nil {
			log.Printf("Reindexing error: %v", err)
		} else {
			log.Println("Reindexing completed successfully")
		}
	}()

	c.JSON(http.StatusAccepted, dto.Response{
		Status:  http.StatusAccepted,
		Message: "Reindexing request received",
		Data: map[string]interface{}{
			"message": "The reindexing process has been triggered. If a process was already running, this request will be ignored to prevent duplicates.",
		},
	})
}
