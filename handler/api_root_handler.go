package handler

import (
	"net/http"

	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/gin-gonic/gin"
)

type ApiRootHandler struct {
	service service.IApiRootService
}

func NewApiRootHandler(ars service.IApiRootService) *ApiRootHandler {
	return &ApiRootHandler{
		service: ars,
	}
}

func (h *ApiRootHandler) GetV1(c *gin.Context) {
	apiRoot, err := h.service.GetV1()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load api index",
		})
		return
	}

	c.JSON(http.StatusOK, apiRoot)
}