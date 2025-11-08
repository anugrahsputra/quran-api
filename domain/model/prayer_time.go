package model

import (
	"strconv"

	"github.com/anugrahsputra/quran-api/domain/dto"
)

type PrayerTime struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Data   PrayerTimeData `json:"data"`
}

type PrayerTimeData struct {
	Timings Timings `json:"timings"`
	Date    Date    `json:"date"`
	Meta    Meta    `json:"meta"`
}

type Timings struct {
	Fajr       string `json:"Fajr"`
	Sunrise    string `json:"Sunrise"`
	Dhuhr      string `json:"Dhuhr"`
	Asr        string `json:"Asr"`
	Sunset     string `json:"Sunset"`
	Maghrib    string `json:"Maghrib"`
	Isha       string `json:"Isha"`
	Imsak      string `json:"Imsak"`
	Midnight   string `json:"Midnight"`
	Firstthird string `json:"Firstthird"`
	Lastthird  string `json:"Lastthird"`
}

type Date struct {
	Readable  string    `json:"readable"`
	Timestamp string    `json:"timestamp"`
	Hijri     Hijri     `json:"hijri"`
	Gregorian Gregorian `json:"gregorian"`
}

type Hijri struct {
	Date             string       `json:"date"`
	Format           string       `json:"format"`
	Day              string       `json:"day"`
	Weekday          HijriWeekday `json:"weekday"`
	Month            HijriMonth   `json:"month"`
	Year             string       `json:"year"`
	Designation      Designation  `json:"designation"`
	Holidays         []any        `json:"holidays"`
	AdjustedHolidays []string     `json:"adjustedHolidays"`
	Method           string       `json:"method"`
}

type HijriWeekday struct {
	En string `json:"en"`
	Ar string `json:"ar"`
}

type HijriMonth struct {
	Number int    `json:"number"`
	En     string `json:"en"`
	Ar     string `json:"ar"`
	Days   int    `json:"days"`
}

type Designation struct {
	Abbreviated string `json:"abbreviated"`
	Expanded    string `json:"expanded"`
}

type Gregorian struct {
	Date          string           `json:"date"`
	Format        string           `json:"format"`
	Day           string           `json:"day"`
	Weekday       GregorianWeekday `json:"weekday"`
	Month         GregorianMonth   `json:"month"`
	Year          string           `json:"year"`
	Designation   Designation      `json:"designation"`
	LunarSighting bool             `json:"lunarSighting"`
}

type GregorianWeekday struct {
	En string `json:"en"`
}

type GregorianMonth struct {
	Number int    `json:"number"`
	En     string `json:"en"`
}

type Meta struct {
	Latitude                 float32 `json:"latitude"`
	Longitude                float32 `json:"longitude"`
	Timezone                 string  `json:"timezone"`
	Method                   Method  `json:"method"`
	LatitudeAdjustmentMethod string  `json:"latitudeAdjustmentMethod"`
	MidnightMode             string  `json:"midnightMode"`
	School                   string  `json:"school"`
	Offset                   Offset  `json:"offset"`
}

type Method struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Params   Params   `json:"params"`
	Location Location `json:"location"`
}

type Params struct {
	Fajr int `json:"Fajr"`
	Isha int `json:"Isha"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Offset struct {
	Imsak    int `json:"Imsak"`
	Fajr     int `json:"Fajr"`
	Sunrise  int `json:"Sunrise"`
	Dhuhr    int `json:"Dhuhr"`
	Asr      int `json:"Asr"`
	Maghrib  int `json:"Maghrib"`
	Sunset   int `json:"Sunset"`
	Isha     int `json:"Isha"`
	Midnight int `json:"Midnight"`
}

func (p *PrayerTime) ToDTO(city string) dto.PrayerTimeResp {
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
