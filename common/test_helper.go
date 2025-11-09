package common

import (
	"context"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/domain/model"
	"github.com/stretchr/testify/mock"
)

// MockSurahService is a mock implementation of ISurahService
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

// MockQuranRepository is a mock implementation of IQuranRepository
type MockQuranRepository struct {
	mock.Mock
}

func (m *MockQuranRepository) GetListSurah(ctx context.Context) ([]model.Surah, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Surah), args.Error(1)
}

func (m *MockQuranRepository) GetSurahDetail(ctx context.Context, id int, start int, pageLimit int) (model.DetailSurahApi, error) {
	args := m.Called(ctx, id, start, pageLimit)
	return args.Get(0).(model.DetailSurahApi), args.Error(1)
}
