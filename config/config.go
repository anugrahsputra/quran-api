package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anugrahsputra/go-quran-api/utils/helper"
	"github.com/op/go-logging"
)

type Config struct {
	Port            string
	SearchIndexPath string
	ExternalUrl     ExternalUrl
}

type ExternalUrl struct {
	KemenagApi    string
	PrayerTimeApi string
}

func LoadConfig() *Config {
	return &Config{
		Port:            helper.GetEnv("PORT", "8080"),
		SearchIndexPath: helper.GetEnv("SEARCH_INDEX_PATH", "quran.bleve"),
		ExternalUrl: ExternalUrl{
			KemenagApi:    helper.GetEnv("KEMENAG_API", "https://web-api.qurankemenag.net"),
			PrayerTimeApi: helper.GetEnv("PRAYER_TIME_API", "https://api.aladhan.com/v1"),
		},
	}
}

func (c *Config) FetchFromKemenag(ctx context.Context, httpClient *http.Client, path string, v any) error {
	url := fmt.Sprintf("%s/%s", c.ExternalUrl.KemenagApi, path)
	logger := logging.MustGetLogger("repository")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Origin", "https://quran.kemenag.go.id")
	req.Header.Set("Accept", "application/json")

	start := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch Kemenag: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	logger.Infof("Fetched from %s in %v", path, duration)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Kemenag responded with: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	return json.Unmarshal(body, v)
}
