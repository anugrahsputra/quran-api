package service

import (
	"context"
	"errors"
	"testing"

	"github.com/anugrahsputra/quran-api/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuranRepository is a mock implementation of IQuranRepository

type MockQuranRepository struct {
	mock.Mock
}

func (m *MockQuranRepository) GetListSurah(ctx context.Context) ([]model.Surah, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Surah), args.Error(1)
}

func TestSurahService_GetListSurah_Success(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(MockQuranRepository)

	// Create a sample surah list
	surahs := []model.Surah{
		{
			ID:    1,
			Latin: "Al-Fatihah",
		},
	}

	// Expect a call to GetListSurah and return the sample surah list
	mockRepo.On("GetListSurah", context.Background()).Return(surahs, nil)

	// Create a new surahService with the mock repository
	service := NewSurahService(mockRepo)

	// Call the method to be tested
	surahsResp, err := service.GetListSurah(context.Background())

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, surahsResp)
	assert.Len(t, surahsResp, 1)
	assert.Equal(t, 1, surahsResp[0].ID)
	assert.Equal(t, "Al-Fatihah", surahsResp[0].Latin)

	// Assert that the expected methods were called
	mockRepo.AssertExpectations(t)
}

func TestSurahService_GetListSurah_Error(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(MockQuranRepository)

	// Expect a call to GetListSurah and return an error
	mockRepo.On("GetListSurah", context.Background()).Return(([]model.Surah)(nil), errors.New("repository error"))

	// Create a new surahService with the mock repository
	service := NewSurahService(mockRepo)

	// Call the method to be tested
	surahsResp, err := service.GetListSurah(context.Background())

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, surahsResp)
	assert.Equal(t, "repository error", err.Error())

	// Assert that the expected methods were called
	mockRepo.AssertExpectations(t)
}
