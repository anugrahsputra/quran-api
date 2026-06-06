package dto

import (
	"github.com/anugrahsputra/go-quran-api/internal/domain"
)

type SearchResponse struct {
	Code    int                   `json:"code"`
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Meta    Meta                  `json:"meta"`
	Data    []domain.SearchedAyah `json:"data"`
}
