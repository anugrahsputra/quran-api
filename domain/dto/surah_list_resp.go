package dto

import "github.com/anugrahsputra/quran-api/domain/model"

type SurahListResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []model.Surah
}
