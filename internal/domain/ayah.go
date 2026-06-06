package domain

import (
	"context"
	"time"
)

type AyahResponse struct {
	Data Ayah `json:"data"`
}

type Ayah struct {
	ID          int
	SurahID     int
	Ayah        int
	Page        int
	QuarterHizb float32
	Juz         int
	Manzil      int
	Arabic      string
	Kitabah     string
	Latin       string
	ArabicWords []string
	Translation string
	Footnotes   *string
	UpdatedAt   time.Time
	Surah       Surah
	Tafsir      Tafsir
}

type Tafsir struct {
	Wajiz              string
	Tahlili            string
	IntroSurah         string
	OutroSurah         string
	MunasabahPrevSurah string
	MunasabahPrevTheme string
	ThemeGroup         string
	Kosakata           string
	SababNuzul         string
	Conclusion         string
}

type AyahRepository interface {
	GetAyah(ctx context.Context, id int) (Ayah, error)
}
