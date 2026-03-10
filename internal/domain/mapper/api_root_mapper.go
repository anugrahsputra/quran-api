package mapper

import (
	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/anugrahsputra/go-quran-api/internal/domain/model"
)

func ToApiRootDTO(ar *model.ApiRoot) dto.ApiRootDTO {
	paths := make(map[string]dto.ApiLinkDTO, len(ar.Paths))
	for k, v := range ar.Paths {
		paths[k] = dto.ApiLinkDTO{
			Method: v.Method,
			Path:   v.Path,
		}
	}

	return dto.ApiRootDTO{
		Version: ar.Version,
		Paths:   paths,
	}
}
