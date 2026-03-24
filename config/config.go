package config

import (
	"github.com/anugrahsputra/go-quran-api/utils/helper"
)

type Config struct {
	Port            string
	SearchIndexPath string
	ExternalUrl     ExternalUrl
	Redis           RedisConfig
}

type ExternalUrl struct {
	KemenagApi    string
	PrayerTimeApi string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       string
}

func LoadConfig() *Config {
	return &Config{
		Port:            helper.GetEnv("PORT", "8080"),
		SearchIndexPath: helper.GetEnv("SEARCH_INDEX_PATH", "quran.bleve"),
		ExternalUrl: ExternalUrl{
			KemenagApi:    helper.GetEnv("KEMENAG_API", "https://web-api.qurankemenag.net"),
			PrayerTimeApi: helper.GetEnv("PRAYER_TIME_API", "https://api.aladhan.com/v1"),
		},
		Redis: RedisConfig{
			Host:     helper.GetEnv("REDIS_HOST", "localhost"),
			Port:     helper.GetEnv("REDIS_PORT", "6379"),
			Password: helper.GetEnv("REDIS_PASSWORD", ""),
			DB:       helper.GetEnv("REDIS_DB", "0"),
		},
	}
}
