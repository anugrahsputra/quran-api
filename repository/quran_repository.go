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
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("repository")

type IQuranRepository interface {
	GetListSurah(ctx context.Context) ([]model.Surah, error)
}

type quranRepository struct {
	kemenagApi string
	httpClient *http.Client
}

func NewQuranRepository(cfg *config.Config) IQuranRepository {
	return &quranRepository{
		kemenagApi: cfg.ExternalUrl.KemenagApi,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *quranRepository) GetListSurah(ctx context.Context) ([]model.Surah, error) {
	url := fmt.Sprintf("%s/quran-surah", r.kemenagApi)

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Origin", "https://quran.kemenag.go.id")
	req.Header.Set("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		logger.Errorf("Failed to fetch Kemenag: %v", err)
		return nil, fmt.Errorf("failed to fetch Kemenag: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	logger.Infof("Fetched surahs from Kemenag in %s", duration)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch Kemenag: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result model.SurahList
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	logger.Infof("Successfully fetched surahs - Count: %d, Duration: %dms", len(result.Data), duration.Milliseconds())

	return result.Data, nil
}
