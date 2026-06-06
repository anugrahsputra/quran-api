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

type surahRepository struct {
	kemenagApi string
	surfClient *surf.Client
}

func NewSurahRepository(cfg *config.Config) domain.SurahRepository {
	client := surf.NewClient().
		Builder().
		Impersonate().
		RandomOS().
		Android().
		Chrome().
		Timeout(10 * time.Second).
		Build().
		Unwrap()

	return &surahRepository{
		kemenagApi: cfg.ExternalUrl.KemenagApi,
		surfClient: client,
	}
}

func (r *surahRepository) GetListSurah(ctx context.Context) ([]domain.Surah, error) {
	var result domain.SurahResponse

	if err := r.fetchFromKemenag(ctx, "quran-surah", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (r *surahRepository) GetSurahDetail(ctx context.Context, id int, start int, pageLimit int) ([]domain.DetailSurah, error) {
	var result domain.DetailSurahResponse
	surahDetail := fmt.Sprintf(common.DETAIL_SURAH, id, start, pageLimit)
	if err := r.fetchFromKemenag(ctx, surahDetail, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (r *surahRepository) fetchFromKemenag(ctx context.Context, path string, v any) error {
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
