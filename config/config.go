package config

import (
	"github.com/anugrahsputra/go-quran-api/utils/helper"
)

type Config struct {
	Port        string
	ExternalUrl ExternalUrl
}

type ExternalUrl struct {
	KemenagApi    string
	PrayerTimeApi string
}

func LoadConfig() *Config {
	return &Config{
		Port: helper.GetEnv("PORT", "8080"),
		ExternalUrl: ExternalUrl{
			KemenagApi:    helper.GetEnv("KEMENAG_API", "https://web-api.qurankemenag.net"),
			PrayerTimeApi: helper.GetEnv("PRAYER_TIME_API", "https://api.aladhan.com/v1"),
		},
	}
}
