package dto

type SurahDetailResp struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Meta    Meta            `json:"meta"`
	Data    SurahDetailData `json:"data"`
}

type SurahDetailData struct {
	SurahID         int     `json:"surah_id"`
	Arabic          string  `json:"arabic"`
	Latin           string  `json:"latin"`
	Translation     string  `json:"translation"`
	Transliteration string  `json:"transliteration"`
	Location        string  `json:"location"`
	Audio           string  `json:"audio"`
	Verses          []Verse `json:"verses"`
}

type AudioResp struct {
	Primary   string   `json:"primary"`
	Secondary []string `json:"secondary"`
}

type Verse struct {
	Id          int     `json:"id"`
	Ayah        int     `json:"ayah"`
	Page        int     `json:"page"`
	QuarterHizb float32 `json:"quarter_hizb"`
	Juz         int     `json:"juz"`
	Manzil      int     `json:"manzil"`
	Arabic      string  `json:"arabic"`
	Kitabah     string  `json:"kitabah"`
	Latin       string  `json:"latin"`
	Translation string  `json:"translation"`
	Audio       string  `json:"audio"`
}
