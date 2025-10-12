package model

import (
	"time"

	"github.com/anugrahsputra/quran-api/domain/dto"
)

type SurahList struct {
	Data []Surah `json:"data"`
}

type Surah struct {
	ID              int       `json:"id"`
	Arabic          string    `json:"arabic"`
	Latin           string    `json:"latin"`
	Transliteration string    `json:"transliteration"`
	Translation     string    `json:"translation"`
	NumAyah         int       `json:"num_ayah"`
	Page            int       `json:"page"`
	Location        string    `json:"location"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (surah *Surah) ToDTO() dto.SurahResp {
	return dto.SurahResp{
		ID:              surah.ID,
		Arabic:          surah.Arabic,
		Latin:           surah.Latin,
		Transliteration: surah.Transliteration,
		Translation:     surah.Translation,
		NumAyah:         surah.NumAyah,
		Page:            surah.Page,
		Location:        surah.Location,
		UpdatedAt:       surah.UpdatedAt,
	}
}
