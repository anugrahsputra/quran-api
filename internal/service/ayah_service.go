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

type AyahService interface {
	GetAyah(ctx context.Context, id int) (dto.DetailAyahResp, error)
}

type ayahService struct {
	repo        domain.AyahRepository
	redisClient *redis.Client
}

func NewAyahService(r domain.AyahRepository, rc *redis.Client) AyahService {
	return &ayahService{
		repo:        r,
		redisClient: rc,
	}
}

func (s *ayahService) GetAyah(ctx context.Context, id int) (dto.DetailAyahResp, error) {
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

	detailAyah, err := s.repo.GetAyah(ctx, id)
	if err != nil {
		return dto.DetailAyahResp{}, err
	}

	response := mapper.ToAyahDTO(&detailAyah)

	if s.redisClient != nil {
		data, _ := json.Marshal(response)
		s.redisClient.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return response, nil
}
