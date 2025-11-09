package model

type Ayah struct {
	SurahNumber int    `json:"surah_number"`
	AyahNumber  int    `json:"ayah_number"`
	Text        string `json:"text"`
	Latin       string `json:"latin"`
	Translation string `json:"translation"`
}
