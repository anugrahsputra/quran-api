package service

import (
	"context"
	"fmt"

	"github.com/anugrahsputra/quran-api/domain/dto"
	"github.com/anugrahsputra/quran-api/repository"
)

type ISurahService interface {
	GetListSurah(ctx context.Context) ([]dto.SurahResp, error)
	GetSurahDetail(ctx context.Context, id int, start int, pageLimit int) (dto.SurahDetailResp, error)
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
	surahs, err := s.repository.GetListSurah(ctx)
	if err != nil {
		return nil, err
	}

	var surahsResp []dto.SurahResp
	for _, surah := range surahs {
		surahsResp = append(surahsResp, surah.ToDTO())
	}

	return surahsResp, nil
}

func (s *surahService) GetSurahDetail(ctx context.Context, id int, start int, pageLimit int) (dto.SurahDetailResp, error) {
	surahApi, err := s.repository.GetSurahDetail(ctx, id, start, pageLimit)
	if err != nil {
		return dto.SurahDetailResp{}, err
	}

	response := surahApi.ToDetailDTO()
	fmt.Println(response)
	return response, nil
}
