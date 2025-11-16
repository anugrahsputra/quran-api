package service

import (
	"context"
	"fmt"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/domain/mapper"
	"github.com/anugrahsputra/go-quran-api/repository"
	"github.com/op/go-logging"
	"go.uber.org/zap"
)

var logger = logging.MustGetLogger("service")

type IQuranService interface {
	GetListSurah(ctx context.Context) ([]dto.SurahResp, error)
	GetSurahDetail(ctx context.Context, id int, page int, limit int) (dto.SurahDetailData, int, int, error)
	GetDetailAyah(ctx context.Context, id int) (dto.DetailAyahResp, error)
}

type quranService struct {
	repository repository.IQuranRepository
}

func NewQuranService(r repository.IQuranRepository) IQuranService {
	return &quranService{
		repository: r,
	}
}

func (s *quranService) GetListSurah(ctx context.Context) ([]dto.SurahResp, error) {
	surahs, err := s.repository.GetListSurah(ctx)
	if err != nil {
		return nil, err
	}

	var surahsResp []dto.SurahResp
	for _, surah := range surahs {
		surahsResp = append(surahsResp, mapper.ToSurahDTO(&surah))
	}

	return surahsResp, nil
}

func (s *quranService) GetSurahDetail(ctx context.Context, id int, page int, limit int) (dto.SurahDetailData, int, int, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Convert page to start (offset) for the repository
	start := (page - 1) * limit

	// First, get the surah list to get the total number of ayahs
	surahs, err := s.repository.GetListSurah(ctx)
	if err != nil {
		return dto.SurahDetailData{}, 0, 0, err
	}

	// Find the surah to get NumAyah (total verses)
	var totalVerses int
	for _, surah := range surahs {
		if surah.ID == id {
			totalVerses = surah.NumAyah
			break
		}
	}

	if totalVerses == 0 {
		return dto.SurahDetailData{}, 0, 0, fmt.Errorf("surah with id %d not found", id)
	}

	// Fetch the verses for the requested page
	surahApi, err := s.repository.GetSurahDetail(ctx, id, start, limit)
	if err != nil {
		return dto.SurahDetailData{}, 0, 0, err
	}

	if len(surahApi.Data) == 0 {
		return dto.SurahDetailData{}, totalVerses, 0, nil
	}

	surahDTO := mapper.ToSurahDTO(&surahApi.Data[0].Surah)
	verses := make([]dto.Verse, len(surahApi.Data))
	for i, verse := range surahApi.Data {
		verses[i] = mapper.ToVerseDTO(&verse)
	}

	// Calculate total pages
	totalPages := (totalVerses + limit - 1) / limit // Ceiling division
	if totalPages == 0 && totalVerses > 0 {
		totalPages = 1
	}

	response := dto.SurahDetailData{
		SurahID:         surahDTO.ID,
		Arabic:          surahDTO.Arabic,
		Latin:           surahDTO.Latin,
		Translation:     surahDTO.Translation,
		Transliteration: surahDTO.Transliteration,
		Location:        surahDTO.Location,
		Audio:           fmt.Sprintf(mapper.SURAH_AUDIO_URL, surahDTO.ID),
		Verses:          verses,
	}

	logger.Info("Fetched surah detail", zap.Int("id", id), zap.Int("page", page), zap.Int("limit", limit), zap.Int("total", totalVerses))
	return response, totalVerses, totalPages, nil
}

func (s *quranService) GetDetailAyah(ctx context.Context, id int) (dto.DetailAyahResp, error) {

	detailAyah, err := s.repository.GetDetailAyah(ctx, id)
	if err != nil {
		return dto.DetailAyahResp{}, err
	}

	response := mapper.ToDetailAyahDTO(&detailAyah)
	return response, nil
}
