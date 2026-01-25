package model

type ApiRoot struct {
	Version string  `json:"version"`
	Paths   ApiPath `json:"paths"`
}

type ApiPath struct {
	ListSurah   string `json:"list_surah"`
	DetailSurah string `json:"detail_surah"`
	Ayah        string `json:"ayah"`
	PrayerTime  string `json:"prayer_time"`
	TopicSearch string `json:"topic_search"`
}
