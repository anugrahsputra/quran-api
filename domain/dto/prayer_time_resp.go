package dto

type PrayerTimeResp struct {
	Date     Date     `json:"date"`
	Location Location `json:"location"`
	Timings  Timings  `json:"timings"`
}

type Timings struct {
	Imsak   string `json:"Imsak"`
	Fajr    string `json:"Fajr"`
	Sunrise string `json:"Sunrise"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

type Date struct {
	Gregorian string `json:"gregorian"`
	Hijri     string `json:"hijri"`
	Weekday   string `json:"weekday"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type Location struct {
	City     string `json:"city"`
	Timezone string `json:"timezone"`
	Method   string `json:"method"`
}
