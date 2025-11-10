package service

import (
	"context"
	"errors"
	"testing"

	"github.com/anugrahsputra/go-quran-api/common"
	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestSurahService_GetListSurah_Success(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(common.MockQuranRepository)

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
	mockRepo := new(common.MockQuranRepository)

	// Expect a call to GetListSurah and return an error
	mockRepo.On("GetListSurah", context.Background()).Return(([]model.Surah)(nil), errors.New("repository error"))

	// Create a new surahService with the mock repository
	service := NewSurahService(mockRepo)

	// Call the method to be tested
	surahsResp, err := service.GetListSurah(context.Background())

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, surahsResp)
	// Assert that the expected methods were called
	mockRepo.AssertExpectations(t)
}

func TestSurahService_GetSurahDetail_Success(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(common.MockQuranRepository)

	// Create a sample surah detail
	detailSurah := model.DetailSurahApi{
		Data: []model.DetailSurah{
			{
				ID: 1,
				Surah: model.Surah{
					ID: 1,
				},
			},
		},
	}

	// Expect a call to GetSurahDetail and return the sample surah detail
	mockRepo.On("GetSurahDetail", context.Background(), 1, 0, 10).Return(detailSurah, nil)

	// Create a new surahService with the mock repository
	service := NewSurahService(mockRepo)

	// Call the method to be tested
	detailSurahResp, err := service.GetSurahDetail(context.Background(), 1, 0, 10)

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, detailSurahResp)
	assert.Equal(t, 1, detailSurahResp.SurahID)

	// Assert that the expected methods were called
	mockRepo.AssertExpectations(t)
}

func TestSurahService_GetSurahDetail_Error(t *testing.T) {
	// Create a new mock repository
	mockRepo := new(common.MockQuranRepository)

	// Expect a call to GetSurahDetail and return an error
	mockRepo.On("GetSurahDetail", context.Background(), 1, 0, 10).Return(model.DetailSurahApi{}, errors.New("repository error"))

	// Create a new surahService with the mock repository
	service := NewSurahService(mockRepo)

	// Call the method to be tested
	detailSurahResp, err := service.GetSurahDetail(context.Background(), 1, 0, 10)

	// Assert the results
	assert.Error(t, err)
	assert.Equal(t, dto.SurahDetailResp{}, detailSurahResp)
	assert.Equal(t, "repository error", err.Error())

	// Assert that the expected methods were called
	mockRepo.AssertExpectations(t)
}
