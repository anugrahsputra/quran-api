package domain

import (
	"context"
	"time"
)

type SurahResponse struct {
	Data []Surah `json:"data"`
}

type DetailSurahResponse struct {
	Data []DetailSurah `json:"data"`
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
	Audio           string    `json:"audio"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DetailSurah struct {
	ID          int       `json:"id"`
	SurahID     int       `json:"surah_id"`
	Ayah        int       `json:"ayah"`
	Page        int       `json:"page"`
	QuarterHizb float32   `json:"quarter_hizb"`
	Juz         int       `json:"juz"`
	Manzil      int       `json:"manzil"`
	Arabic      string    `json:"arabic"`
	Kitabah     string    `json:"kitabah"`
	Latin       string    `json:"latin"`
	ArabicWords []string  `json:"arabic_words"`
	Translation string    `json:"translation"`
	Footnotes   *string   `json:"footnotes"`
	UpdatedAt   time.Time `json:"updated_at"`
	Surah       Surah     `json:"surah"`
}

type SurahRepository interface {
	GetListSurah(ctx context.Context) ([]Surah, error)
	GetSurahDetail(ctx context.Context, id int, start int, pageLimit int) ([]DetailSurah, error)
}
