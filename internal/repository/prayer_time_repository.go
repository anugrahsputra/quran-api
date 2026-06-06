package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
	"github.com/enetx/g"
	"github.com/enetx/surf"
)

const PRAYER_TIME_API = "timingsByAddress?address=%s&timezonestring=%s"

// type IPrayerTimeRepository interface {
// 	GetPrayerTime(ctx context.Context, city string, timezone string) (domain.PrayerTime, error)
// }

type prayerTimeRepository struct {
	prayerTimeApi string
	surfClient    *surf.Client
}

func NewPrayerTimeRepository(cfg *config.Config) domain.PrayerTimeRepository {
	client := surf.NewClient().
		Builder().
		Impersonate().
		RandomOS().
		Android().
		Chrome().
		Timeout(10 * time.Second).
		Build().
		Unwrap()

	return &prayerTimeRepository{
		prayerTimeApi: cfg.ExternalUrl.PrayerTimeApi,
		surfClient:    client,
	}
}

func (r *prayerTimeRepository) GetPrayerTime(ctx context.Context, city string, timezone string) (domain.PrayerTime, error) {
	var result domain.PrayerTimeResponse

	prayerTime := fmt.Sprintf(PRAYER_TIME_API, city, timezone)
	if err := r.fetchFromPrayerTimeApi(ctx, prayerTime, &result); err != nil {
		return domain.PrayerTime{}, err
	}
	return result.Data, nil
}

func (r *prayerTimeRepository) fetchFromPrayerTimeApi(ctx context.Context, path string, v any) error {
	url := fmt.Sprintf("%s/%s", r.prayerTimeApi, path)

	resp := r.surfClient.Get(g.String(url)).
		WithContext(ctx).
		Do()

	if resp.IsErr() {
		return fmt.Errorf("failed to fetch prayer time data: %w", resp.Err())
	}

	result := resp.Ok()

	if result.StatusCode != http.StatusOK {
		return fmt.Errorf("prayer time API responded with: %v", result.StatusCode)
	}

	bodyResult := result.Body.Bytes()
	if bodyResult.IsErr() {
		return fmt.Errorf("failed to read body: %w", bodyResult.Err())
	}

	return json.Unmarshal(bodyResult.Ok(), v)
}
