package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/domain/model"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("repository")

const DETAIL_SURAH = "quran-ayah?surah=%d&start=%d&limit=%d"
const DETAIL_AYAH = "quran-tafsir/%d"

type IQuranRepository interface {
	GetListSurah(ctx context.Context) ([]model.Surah, error)
	GetSurahDetail(ctx context.Context, id int, start int, pageLimit int) (model.DetailSurahApi, error)
	GetDetailAyah(ctx context.Context, id int) (model.TafsirData, error)
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
	var result model.SurahList
	if err := r.fetchFromKemenag(ctx, "quran-surah", &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (r *quranRepository) GetSurahDetail(ctx context.Context, id int, start int, pageLimit int) (model.DetailSurahApi, error) {
	var result model.DetailSurahApi

	quranDetail := fmt.Sprintf(DETAIL_SURAH, id, start, pageLimit)
	if err := r.fetchFromKemenag(ctx, quranDetail, &result); err != nil {
		return model.DetailSurahApi{}, err
	}
	return result, nil
}

func (r *quranRepository) GetDetailAyah(ctx context.Context, id int) (model.TafsirData, error) {
	var result model.TafsirApi

	ayahDetail := fmt.Sprintf(DETAIL_AYAH, id)
	if err := r.fetchFromKemenag(ctx, ayahDetail, &result); err != nil {
		return model.TafsirData{}, err
	}
	return result.Data, nil
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
