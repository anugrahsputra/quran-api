package repository

import "github.com/anugrahsputra/go-quran-api/domain/model"

type IApiRootRepository interface {
	GetV1() (*model.ApiRoot, error)
}

type apiRootRepository struct{}

func NewApiRootRepository() IApiRootRepository {
	return &apiRootRepository{}
}

func (r *apiRootRepository) GetV1() (*model.ApiRoot, error) {
	return &model.ApiRoot{
		Version: "v1",
		Paths: map[string]model.ApiLink{
			"list_surah": {
				Method: "GET",
				Path:   "/api/v1/surah",
			},
			"detail_surah": {
				Method: "GET",
				Path:   "/api/v1/surah/:id",
			},
			"ayah": {
				Method: "GET",
				Path:   "/api/v1/ayah/:id",
			},
			"search": {
				Method: "GET",
				Path:   "/api/v1/search?q={query}",
			},
			"prayer_time": {
				Method: "GET",
				Path:   "/api/v1/prayer-time",
			},
		},
	}, nil
}
