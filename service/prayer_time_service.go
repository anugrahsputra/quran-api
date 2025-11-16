package service

import (
	"context"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/domain/mapper"
	"github.com/anugrahsputra/go-quran-api/repository"
)

type IPrayerTimeService interface {
	GetPrayerTime(ctx context.Context, city string, timezone string) (dto.PrayerTimeResp, error)
}

type prayerTimeService struct {
	repository repository.IPrayerTimeRepository
}

func NewPrayerTimeService(r repository.IPrayerTimeRepository) IPrayerTimeService {
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
