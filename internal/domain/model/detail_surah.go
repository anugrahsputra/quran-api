package model

import (
	"time"
)

type DetailSurahApi struct {
	Data []DetailSurah `json:"data"`
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
