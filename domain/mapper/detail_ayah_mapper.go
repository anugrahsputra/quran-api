package mapper

import (
	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/domain/model"
)

// ToDetailAyahDTO converts a model.TafsirData to dto.DetailAyahResp
func ToDetailAyahDTO(da *model.TafsirData) dto.DetailAyahResp {
	return dto.DetailAyahResp{
		ID:          da.ID,
		SurahID:     da.SurahID,
		Ayah:        da.Ayah,
		Page:        da.Page,
		QuarterHizb: da.QuarterHizb,
		Juz:         da.Juz,
		Manzil:      da.Manzil,
		Arabic:      da.Arabic,
		Kitabah:     da.Kitabah,
		Latin:       da.Latin,
		ArabicWords: da.ArabicWords,
		Translation: da.Translation,
		Surah:       ToSurahDTO(&da.Surah),
		Tafsir:      ToTafsirDTO(&da.Tafsir),
	}
}

// ToTafsirDTO converts a model.Tafsir to dto.Tafsir
func ToTafsirDTO(tafsir *model.Tafsir) dto.Tafsir {
	return dto.Tafsir{
		Wajiz:              tafsir.Wajiz,
		Tahlili:            tafsir.Tahlili,
		IntroSurah:         tafsir.IntroSurah,
		OutroSurah:         tafsir.OutroSurah,
		MunasabahPrevSurah: tafsir.MunasabahPrevSurah,
		MunasabahPrevTheme: tafsir.MunasabahPrevTheme,
		ThemeGroup:         tafsir.ThemeGroup,
		Kosakata:           tafsir.Kosakata,
		SababNuzul:         tafsir.SababNuzul,
		Conclusion:         tafsir.Conclusion,
	}
}
