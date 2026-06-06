package service

import (
	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
	"github.com/anugrahsputra/go-quran-api/internal/mapper"
)

type IApiRootService interface {
	GetV1() (dto.ApiRootDTO, error)
}

type apiRootService struct {
	repository domain.ApiRootRepository
}

func NewApiRootService(r domain.ApiRootRepository) IApiRootService {
	return &apiRootService{
		repository: r,
	}
}

func (s *apiRootService) GetV1() (dto.ApiRootDTO, error) {
	apiRoot, err := s.repository.GetV1()
	if err != nil {
		return dto.ApiRootDTO{}, err
	}

	apiRootDto := mapper.ToApiRootDTO(apiRoot)
	return apiRootDto, nil
}
