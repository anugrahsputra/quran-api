package mapper

import (
	"fmt"

	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
)

const SURAH_AUDIO_URL = "https://cdn.islamic.network/quran/audio-surah/128/ar.alafasy/%d.mp3"
const AYAH_AUDIO_URL = "https://cdn.islamic.network/quran/audio/128/ar.alafasy/%d.mp3"

func ToSurahDTO(surah *domain.Surah) dto.SurahResp {
	return dto.SurahResp{
		ID:              surah.ID,
		Arabic:          surah.Arabic,
		Latin:           surah.Latin,
		Transliteration: surah.Transliteration,
		Translation:     surah.Translation,
		NumAyah:         surah.NumAyah,
		Page:            surah.Page,
		Audio:           fmt.Sprintf(SURAH_AUDIO_URL, surah.ID),
		Location:        surah.Location,
		UpdatedAt:       surah.UpdatedAt,
	}
}

func ToVerseDTO(detailSurah *domain.DetailSurah) dto.Verse {
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
