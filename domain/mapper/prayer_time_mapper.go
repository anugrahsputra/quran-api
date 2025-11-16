package mapper

import (
	"strconv"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/domain/model"
)

// ToPrayerTimeDTO converts a model.PrayerTime to dto.PrayerTimeResp
func ToPrayerTimeDTO(p *model.PrayerTime, city string) dto.PrayerTimeResp {
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
