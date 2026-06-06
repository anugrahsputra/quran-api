package mapper

import (
	"strconv"

	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
)

func ToPrayerTimeDTO(p *domain.PrayerTime, city string) dto.PrayerTimeResp {
	if p.Timings.Fajr == "" {
		return dto.PrayerTimeResp{}
	}

	date := p.Date
	hijri := date.Hijri
	greg := date.Gregorian
	meta := p.Meta
	timings := p.Timings

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
