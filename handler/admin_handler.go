package handler

import (
	"log"
	"net/http"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	searchAyahService service.QuranSearchService
}

func NewAdminHandler(searchAyahService service.QuranSearchService) *AdminHandler {
	return &AdminHandler{
		searchAyahService: searchAyahService,
	}
}

func (h *AdminHandler) Reindex(c *gin.Context) {
	// Attempt to start indexing in a non-blocking way to check status
	// But since the service's IndexQuran is now protected, we should ideally
	// check the state first or handle the error from the service.
	
	// We'll run the indexing in background, but the service will prevent duplicates.
	// To give immediate feedback, we'll try to trigger a "dry run" or use the service's lock.
	
	// Better approach: Let's make IndexQuran return quickly if locked, 
	// or provide an IsIndexing() method.
	// For now, we'll just handle the background log.
	
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
