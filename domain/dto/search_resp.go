package dto

import "github.com/anugrahsputra/go-quran-api/domain/model"

type SearchResponse struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Meta    Meta         `json:"meta"`
	Data    []model.Ayah `json:"data"`
}
