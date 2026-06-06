package service

import (
	"context"

	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
	"github.com/anugrahsputra/go-quran-api/internal/mapper"
)

type IPrayerTimeService interface {
	GetPrayerTime(ctx context.Context, city string, timezone string) (dto.PrayerTimeResp, error)
}

type prayerTimeService struct {
	repository domain.PrayerTimeRepository
}

func NewPrayerTimeService(r domain.PrayerTimeRepository) IPrayerTimeService {
	return &prayerTimeService{
		repository: r,
	}
}

func (s *prayerTimeService) GetPrayerTime(ctx context.Context, city string, timezone string) (dto.PrayerTimeResp, error) {
	prayerTime, err := s.repository.GetPrayerTime(ctx, city, timezone)
	if err != nil {
		return dto.PrayerTimeResp{}, err
	}

	response := mapper.ToPrayerTimeDTO(&prayerTime, city)
	return response, nil
}
