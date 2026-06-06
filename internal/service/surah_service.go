package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
	"github.com/anugrahsputra/go-quran-api/internal/mapper"
	"github.com/redis/go-redis/v9"
)

type SurahService interface {
	GetListSurah(ctx context.Context) ([]dto.SurahResp, error)
	GetSurahDetail(ctx context.Context, id int, page int, limit int) (dto.SurahDetailData, int, int, error)
}

type surahService struct {
	repo domain.SurahRepository
	rc   *redis.Client
}

func NewSurahService(r domain.SurahRepository, rc *redis.Client) SurahService {
	return &surahService{
		repo: r,
		rc:   rc,
	}
}

func (s *surahService) GetListSurah(ctx context.Context) ([]dto.SurahResp, error) {
	cacheKey := "quran:surah:list"

	if s.rc != nil {
		val, err := s.rc.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedData []dto.SurahResp
			if err := json.Unmarshal([]byte(val), &cachedData); err == nil {
				return cachedData, nil
			}
		}
	}

	surahs, err := s.repo.GetListSurah(ctx)
	if err != nil {
		return nil, err
	}

	var surahsResp []dto.SurahResp
	for _, surah := range surahs {
		surahsResp = append(surahsResp, mapper.ToSurahDTO(&surah))
	}

	if s.rc != nil {
		data, _ := json.Marshal(surahsResp)
		s.rc.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return surahsResp, nil
}

func (s *surahService) GetSurahDetail(ctx context.Context, id int, page int, limit int) (dto.SurahDetailData, int, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	cacheKey := fmt.Sprintf("quran:surah:detail:%d:%d:%d", id, page, limit)

	type cacheData struct {
		Response    dto.SurahDetailData `json:"response"`
		TotalVerses int                 `json:"total_verses"`
		TotalPages  int                 `json:"total_pages"`
	}

	if s.rc != nil {
		val, err := s.rc.Get(ctx, cacheKey).Result()
		if err == nil {
			var cached cacheData
			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				return cached.Response, cached.TotalVerses, cached.TotalPages, nil
			}
		}
	}

	start := (page - 1) * limit

	surahs, err := s.repo.GetListSurah(ctx)
	if err != nil {
		return dto.SurahDetailData{}, 0, 0, err
	}

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

	surahApi, err := s.repo.GetSurahDetail(ctx, id, start, limit)
	if err != nil {
		return dto.SurahDetailData{}, 0, 0, err
	}

	if len(surahApi) == 0 {
		return dto.SurahDetailData{}, totalVerses, 0, nil
	}

	surahDTO := mapper.ToSurahDTO(&surahApi[0].Surah)
	verses := make([]dto.Verse, len(surahApi))
	for i, verse := range surahApi {
		verses[i] = mapper.ToVerseDTO(&verse)
	}

	totalPages := (totalVerses + limit - 1) / limit
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

	if s.rc != nil {
		cached := cacheData{
			Response:    response,
			TotalVerses: totalVerses,
			TotalPages:  totalPages,
		}
		data, _ := json.Marshal(cached)
		s.rc.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return response, totalVerses, totalPages, nil
}
