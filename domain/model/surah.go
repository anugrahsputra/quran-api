package model

import (
	"time"
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

