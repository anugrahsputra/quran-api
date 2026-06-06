package domain

import "github.com/blevesearch/bleve/v2"

type SearchedAyah struct {
	SurahNumber int    `json:"surah_number"`
	AyahNumber  int    `json:"ayah_number"`
	Text        string `json:"text"`
	Latin       string `json:"latin"`
	Translation string `json:"translation"`
	Tafsir      string `json:"tafsir"`
	Topic       string `json:"topic"`
}

type QuranSearchRepository interface {
	Index(ayahs []SearchedAyah) error
	Search(q string, page, limit int) (*bleve.SearchResult, error)
	GetDocCount() (uint64, error)
	IsHealthy() bool
}
