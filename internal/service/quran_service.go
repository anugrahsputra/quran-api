package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/anugrahsputra/go-quran-api/internal/domain/mapper"
	"github.com/anugrahsputra/go-quran-api/internal/repository"
	"github.com/op/go-logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var logger = logging.MustGetLogger("service")

type IQuranService interface {
	GetListSurah(ctx context.Context) ([]dto.SurahResp, error)
	GetSurahDetail(ctx context.Context, id int, page int, limit int) (dto.SurahDetailData, int, int, error)
	GetDetailAyah(ctx context.Context, id int) (dto.DetailAyahResp, error)
}

type quranService struct {
	repository  repository.IQuranRepository
	redisClient *redis.Client
}

func NewQuranService(r repository.IQuranRepository, redisClient *redis.Client) IQuranService {
	return &quranService{
		repository:  r,
		redisClient: redisClient,
	}
}

func (s *quranService) GetListSurah(ctx context.Context) ([]dto.SurahResp, error) {
	cacheKey := "quran:surah:list"

	if s.redisClient != nil {
		val, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedData []dto.SurahResp
			if err := json.Unmarshal([]byte(val), &cachedData); err == nil {
				return cachedData, nil
			}
		}
	}

	surahs, err := s.repository.GetListSurah(ctx)
	if err != nil {
		return nil, err
	}

	var surahsResp []dto.SurahResp
	for _, surah := range surahs {
		surahsResp = append(surahsResp, mapper.ToSurahDTO(&surah))
	}

	if s.redisClient != nil {
		data, _ := json.Marshal(surahsResp)
		s.redisClient.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return surahsResp, nil
}

func (s *quranService) GetSurahDetail(ctx context.Context, id int, page int, limit int) (dto.SurahDetailData, int, int, error) {
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

	if s.redisClient != nil {
		val, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cached cacheData
			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				return cached.Response, cached.TotalVerses, cached.TotalPages, nil
			}
		}
	}

	start := (page - 1) * limit

	surahs, err := s.repository.GetListSurah(ctx)
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

	if s.redisClient != nil {
		cached := cacheData{
			Response:    response,
			TotalVerses: totalVerses,
			TotalPages:  totalPages,
		}
		data, _ := json.Marshal(cached)
		s.redisClient.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	logger.Info("Fetched surah detail", zap.Int("id", id), zap.Int("page", page), zap.Int("limit", limit), zap.Int("total", totalVerses))
	return response, totalVerses, totalPages, nil
}

func (s *quranService) GetDetailAyah(ctx context.Context, id int) (dto.DetailAyahResp, error) {
	cacheKey := fmt.Sprintf("quran:ayah:detail:%d", id)

	if s.redisClient != nil {
		val, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedData dto.DetailAyahResp
			if err := json.Unmarshal([]byte(val), &cachedData); err == nil {
				return cachedData, nil
			}
		}
	}

	detailAyah, err := s.repository.GetDetailAyah(ctx, id)
	if err != nil {
		return dto.DetailAyahResp{}, err
	}

	response := mapper.ToDetailAyahDTO(&detailAyah)

	if s.redisClient != nil {
		data, _ := json.Marshal(response)
		s.redisClient.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return response, nil
}
