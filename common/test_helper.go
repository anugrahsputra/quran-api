package common

import (
	"context"

	"github.com/anugrahsputra/go-quran-api/internal/delivery/dto"
	"github.com/anugrahsputra/go-quran-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSurahService struct {
	mock.Mock
}

func (m *MockSurahService) GetListSurah(ctx context.Context) ([]dto.SurahResp, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SurahResp), args.Error(1)
}

func (m *MockSurahService) GetSurahDetail(ctx context.Context, id int, start int, limit int) (dto.SurahDetailResp, error) {
	args := m.Called(ctx, id, start, limit)
	if args.Get(0) == nil {
		return dto.SurahDetailResp{}, args.Error(1)
	}
	return args.Get(0).(dto.SurahDetailResp), args.Error(1)
}

type MockQuranRepository struct {
	mock.Mock
}

func (m *MockQuranRepository) GetListSurah(ctx context.Context) ([]domain.Surah, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Surah), args.Error(1)
}

func (m *MockQuranRepository) GetSurahDetail(ctx context.Context, id int, start int, pageLimit int) (domain.DetailSurahResponse, error) {
	args := m.Called(ctx, id, start, pageLimit)
	return args.Get(0).(domain.DetailSurahResponse), args.Error(1)
}
