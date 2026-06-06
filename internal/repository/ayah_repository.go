package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anugrahsputra/go-quran-api/common"
	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
	"github.com/enetx/g"
	"github.com/enetx/surf"
)

type ayahRepository struct {
	kemenagApi string
	surfClient *surf.Client
}

func NewAyahRepository(cfg *config.Config) domain.AyahRepository {
	client := surf.NewClient().
		Builder().
		Impersonate().
		RandomOS().
		Android().
		Chrome().
		Timeout(10 * time.Second).
		Build().
		Unwrap()

	return &ayahRepository{
		kemenagApi: cfg.ExternalUrl.KemenagApi,
		surfClient: client,
	}
}

func (r *ayahRepository) GetAyah(ctx context.Context, id int) (domain.Ayah, error) {
	var result domain.AyahResponse

	ayah := fmt.Sprintf(common.DETAIL_AYAH, id)
	if err := r.fetchFromKemenag(ctx, ayah, &result); err != nil {
		return domain.Ayah{}, nil
	}

	return result.Data, nil
}

func (r *ayahRepository) fetchFromKemenag(ctx context.Context, path string, v any) error {
	url := fmt.Sprintf("%s/%s", r.kemenagApi, path)

	resp := r.surfClient.Get(g.String(url)).
		WithContext(ctx).
		AddHeaders(map[string]string{
			"Origin": "https://quran.kemenag.go.id",
			"Accept": "application/json",
		}).
		Do()

	if resp.IsErr() {
		return fmt.Errorf("failed to fetch Kemenag: %w", resp.Err())
	}

	result := resp.Ok()

	if result.StatusCode != http.StatusOK {
		return fmt.Errorf("Kemenag responded with: %v", result.StatusCode)
	}

	bodyResult := result.Body.Bytes()
	if bodyResult.IsErr() {
		return fmt.Errorf("failed to read body: %w", bodyResult.Err())
	}

	return json.Unmarshal(bodyResult.Ok(), v)
}
