package service

import (
	"context"

	"github.com/anugrahsputra/quran-api/domain/dto"
	"github.com/anugrahsputra/quran-api/domain/model"
	"github.com/anugrahsputra/quran-api/repository"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("service")

type ISurahService interface {
	GetListSurah(ctx context.Context) (*dto.SurahListResp, error)
}

type surahService struct {
	repository repository.IQuranRepository
}

func NewSurahService(r repository.IQuranRepository) ISurahService {
	return &surahService{
		repository: r,
	}
}

func (s *surahService) GetListSurah(ctx context.Context) (*dto.SurahListResp, error) {
	logger.Info("Starting GetListSurah service")

	logger.Debug("Fetching surahs from kemenag")
	surahs, err := s.repository.GetListSurah(ctx)
	if err != nil {
		logger.Errorf("Failed to fetch surahs from kemenag: %v", err)
		return &dto.SurahListResp{
			Status:  500,
			Message: "Failed to fetch surahs from kemenag",
		}, err
	}

	logger.Infof("Successfully fetched %d surahs from kemenag", len(surahs))

	var surahResp []model.Surah
	for _, surah := range surahs {
		surahResp = append(surahResp, surah.ToDTO())
	}

	return &dto.SurahListResp{
		Status:  200,
		Message: "Successfully fetched surahs from kemenag",
		Data:    surahResp,
	}, nil

}
