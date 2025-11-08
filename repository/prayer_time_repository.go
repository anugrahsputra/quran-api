package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anugrahsputra/quran-api/config"
	"github.com/anugrahsputra/quran-api/domain/model"
	"github.com/anugrahsputra/quran-api/utils/helper"
	"github.com/patrickmn/go-cache"
)

const PRAYER_TIME_API = "timingsByAddress?address=%s&timezonestring=%s"

type IPrayerTimeRepository interface {
	GetPrayerTime(ctx context.Context, city string, timezone string) (model.PrayerTime, error)
}

type prayerTimeRepository struct {
	prayerTimeApi string
	httpClient    *http.Client
	cache         *cache.Cache
}

func NewPrayerTimeRepository(cfg *config.Config) IPrayerTimeRepository {
	return &prayerTimeRepository{
		prayerTimeApi: cfg.ExternalUrl.PrayerTimeApi,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		cache:         cache.New(1*time.Hour, 10*time.Minute),
	}
}

func (r *prayerTimeRepository) GetPrayerTime(ctx context.Context, city string, timezone string) (model.PrayerTime, error) {
	cacheKey := "prayer_time"

	return helper.GetOrSetCache(r.cache, cacheKey, time.Hour, func() (model.PrayerTime, error) {
		var result model.PrayerTime

		prayerTime := fmt.Sprintf(PRAYER_TIME_API, city, timezone)
		if err := r.fetchFromPrayerTimeApi(ctx, prayerTime, &result); err != nil {
			return model.PrayerTime{}, err
		}
		return result, nil
	})

}

func (r *prayerTimeRepository) fetchFromPrayerTimeApi(ctx context.Context, path string, v any) error {
	url := fmt.Sprintf("%s/%s", r.prayerTimeApi, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	start := time.Now()
	resp, err := r.httpClient.Do(req)
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
