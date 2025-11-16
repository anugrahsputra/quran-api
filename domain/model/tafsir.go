package model

import (
	"time"
)

type TafsirApi struct {
	Data TafsirData `json:"data"`
}

type TafsirData struct {
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
	Tafsir      Tafsir    `json:"tafsir"`
}

type Tafsir struct {
	Wajiz              string `json:"wajiz"`
	Tahlili            string `json:"tahlili"`
	IntroSurah         string `json:"intro_surah"`
	OutroSurah         string `json:"outro_surah"`
	MunasabahPrevSurah string `json:"munasabah_prev_surah"`
	MunasabahPrevTheme string `json:"munasabah_prev_theme"`
	ThemeGroup         string `json:"theme_group"`
	Kosakata           string `json:"kosakata"`
	SababNuzul         string `json:"sabab_nuzul"`
	Conclusion         string `json:"conclusion"`
}
