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
	"github.com/op/go-logging"
	"github.com/patrickmn/go-cache"
)

var logger = logging.MustGetLogger("repository")

type IQuranRepository interface {
	GetListSurah(ctx context.Context) ([]model.Surah, error)
}

type quranRepository struct {
	kemenagApi string
	httpClient *http.Client
	cache      *cache.Cache
}

func NewQuranRepository(cfg *config.Config) IQuranRepository {
	return &quranRepository{
		kemenagApi: cfg.ExternalUrl.KemenagApi,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: cache.New(1*time.Hour, 10*time.Minute),
	}
}

func (r *quranRepository) GetListSurah(ctx context.Context) ([]model.Surah, error) {
	cacheKey := "surah_list"

	return helper.GetOrSetCache(r.cache, cacheKey, time.Hour, func() ([]model.Surah, error) {
		var result model.SurahList
		if err := r.fetchFromKemenag(ctx, "quran-surah", &result); err != nil {
			return nil, err
		}
		return result.Data, nil
	})
}

func (r *quranRepository) fetchFromKemenag(ctx context.Context, path string, v any) error {
	url := fmt.Sprintf("%s/%s", r.kemenagApi, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Origin", "https://quran.kemenag.go.id")
	req.Header.Set("Accept", "application/json")

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
