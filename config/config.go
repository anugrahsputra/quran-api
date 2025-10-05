package config

import (
	"github.com/anugrahsputra/quran-api/utils/helper"
)

type Config struct {
	Port        string
	ExternalUrl ExternalUrl
}

type ExternalUrl struct {
	KemenagApi string
}

func LoadConfig() *Config {
	return &Config{
		Port: helper.GetEnv("PORT", "8080"),
		ExternalUrl: ExternalUrl{
			KemenagApi: helper.GetEnv("KEMENAG_API", "https://web-api.qurankemenag.net"),
		},
	}
}
