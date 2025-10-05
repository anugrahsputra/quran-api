package service

import (
	"context"

	"github.com/anugrahsputra/quran-api/domain/dto"
	"github.com/anugrahsputra/quran-api/repository"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("service")

type ISurahService interface {
	GetListSurah(ctx context.Context) ([]dto.SurahResp, error)
}

type surahService struct {
	repository repository.IQuranRepository
}

func NewSurahService(r repository.IQuranRepository) ISurahService {
	return &surahService{
		repository: r,
	}
}

func (s *surahService) GetListSurah(ctx context.Context) ([]dto.SurahResp, error) {
	surash, err := s.repository.GetListSurah(ctx)
	if err != nil {
		return nil, err
	}

	var surahsResp []dto.SurahResp
	for _, surah := range surash {
		surahsResp = append(surahsResp, surah.ToDTO())
	}

	return surahsResp, nil
}
