package service

import (
	"context"
	"strconv"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/domain/model"
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

	response := s.toPrayerTimeDTO(&prayerTime, city)
	return response, nil
}

func (s *prayerTimeService) toPrayerTimeDTO(p *model.PrayerTime, city string) dto.PrayerTimeResp {
	if p.Data.Timings.Fajr == "" { // basic sanity check
		return dto.PrayerTimeResp{}
	}

	date := p.Data.Date
	hijri := date.Hijri
	greg := date.Gregorian
	meta := p.Data.Meta
	timings := p.Data.Timings

	// convert timestamp string to int64 safely
	var timestamp int64
	if t, err := strconv.ParseInt(date.Timestamp, 10, 64); err == nil {
		timestamp = t
	}

	return dto.PrayerTimeResp{
		Date: dto.Date{
			Gregorian: greg.Date,
			Hijri:     hijri.Day + " " + hijri.Month.En + " " + hijri.Year,
			Weekday:   greg.Weekday.En,
			Timestamp: timestamp,
		},
		Location: dto.Location{
			City:     city,
			Timezone: meta.Timezone,
			Method:   meta.Method.Name,
		},
		Timings: dto.Timings{
			Imsak:   timings.Imsak,
			Fajr:    timings.Fajr,
			Sunrise: timings.Sunrise,
			Dhuhr:   timings.Dhuhr,
			Asr:     timings.Asr,
			Maghrib: timings.Maghrib,
			Isha:    timings.Isha,
		},
	}
}
