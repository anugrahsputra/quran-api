package service

import (
	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/domain/mapper"
	"github.com/anugrahsputra/go-quran-api/repository"
)

type IApiRootService interface {
	GetV1() (dto.ApiRootDTO, error)
}

type apiRootService struct {
	repository repository.IApiRootRepository
}

func NewApiRootService(r repository.IApiRootRepository) IApiRootService {
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
	response := dto.ApiRootDTO{
		Version: apiRoot.Version,
		Paths:   apiRootDto.Paths,
	}

	return response, nil
}
