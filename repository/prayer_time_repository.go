package repository
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/domain/model"
	"github.com/enetx/g"
	"github.com/enetx/surf"
)

const PRAYER_TIME_API = "timingsByAddress?address=%s&timezonestring=%s"

type IPrayerTimeRepository interface {
	GetPrayerTime(ctx context.Context, city string, timezone string) (model.PrayerTime, error)
}

type prayerTimeRepository struct {
	prayerTimeApi string
	surfClient    *surf.Client
}

func NewPrayerTimeRepository(cfg *config.Config) IPrayerTimeRepository {
	client := surf.NewClient().
		Builder().
		Timeout(10 * time.Second).
		Build().
		Unwrap()

	return &prayerTimeRepository{
		prayerTimeApi: cfg.ExternalUrl.PrayerTimeApi,
		surfClient:    client,
	}
}

func (r *prayerTimeRepository) GetPrayerTime(ctx context.Context, city string, timezone string) (model.PrayerTime, error) {
	var result model.PrayerTime

	prayerTime := fmt.Sprintf(PRAYER_TIME_API, city, timezone)
	if err := r.fetchFromPrayerTimeApi(ctx, prayerTime, &result); err != nil {
		return model.PrayerTime{}, err
	}
	return result, nil
}

func (r *prayerTimeRepository) fetchFromPrayerTimeApi(ctx context.Context, path string, v any) error {
	url := fmt.Sprintf("%s/%s", r.prayerTimeApi, path)

	start := time.Now()
	resp := r.surfClient.Get(g.String(url)).
		WithContext(ctx).
		Do()

	if resp.IsErr() {
		return fmt.Errorf("failed to fetch prayer time data: %w", resp.Err())
	}

	result := resp.Ok()
	duration := time.Since(start)
	logger.Infof("Fetched from %s in %v", path, duration)

	if result.StatusCode != http.StatusOK {
		return fmt.Errorf("prayer time API responded with: %v", result.StatusCode)
	}

	bodyResult := result.Body.Bytes()
	if bodyResult.IsErr() {
		return fmt.Errorf("failed to read body: %w", bodyResult.Err())
	}

	return json.Unmarshal(bodyResult.Ok(), v)
}
