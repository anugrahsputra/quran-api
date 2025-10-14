package model

import (
	"fmt"
	"time"

	"github.com/anugrahsputra/quran-api/domain/dto"
)

const SURAH_AUDIO_URL = "https://cdn.islamic.network/quran/audio-surah/128/ar.alafasy/%d.mp3"
const AYAH_AUDIO_URL = "https://cdn.islamic.network/quran/audio/128/ar.alafasy/%d.mp3"

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

func (api *DetailSurahApi) ToDetailDTO() dto.SurahDetailResp {
	if len(api.Data) == 0 {
		return dto.SurahDetailResp{}
	}

	surahDTO := api.Data[0].Surah.ToDTO()
	verses := make([]dto.Verse, len(api.Data))
	for i, verse := range api.Data {
		verses[i] = verse.toVerseDTO()
	}

	return dto.SurahDetailResp{
		SurahID:         surahDTO.ID,
		Arabic:          surahDTO.Arabic,
		Latin:           surahDTO.Latin,
		Translation:     surahDTO.Translation,
		Transliteration: surahDTO.Transliteration,
		Location:        surahDTO.Location,
		Audio:           fmt.Sprintf(SURAH_AUDIO_URL, surahDTO.ID),
		Verses:          verses,
	}
}

func (detailSurah *DetailSurah) toVerseDTO() dto.Verse {
	return dto.Verse{
		Id:          detailSurah.ID,
		Ayah:        detailSurah.Ayah,
		Page:        detailSurah.Page,
		QuarterHizb: detailSurah.QuarterHizb,
		Juz:         detailSurah.Juz,
		Manzil:      detailSurah.Manzil,
		Arabic:      detailSurah.Arabic,
		Kitabah:     detailSurah.Kitabah,
		Latin:       detailSurah.Latin,
		Translation: detailSurah.Translation,
		Audio:       fmt.Sprintf(AYAH_AUDIO_URL, detailSurah.ID),
	}
}
