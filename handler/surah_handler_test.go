package handler

import (
	"time"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anugrahsputra/quran-api/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSurahService is a mock implementation of ISurahService

type MockSurahService struct {
	mock.Mock
}

func (m *MockSurahService) GetListSurah(ctx context.Context) ([]dto.SurahResp, error) {
	args := m.Called(ctx)
	// Handle the case where the first return value is nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.SurahResp), args.Error(1)
}

func TestSurahHandler_GetListSurah_Success(t *testing.T) {
	// Create a new mock service
	mockService := new(MockSurahService)

	// Create a sample surah list
	surahs := []dto.SurahResp{
		{
			ID:    1,
			Latin: "Al-Fatihah",
		},
	}

	// Expect a call to GetListSurah and return the sample surah list
	mockService.On("GetListSurah", mock.Anything).Return(surahs, nil)

	// Create a new surahHandler with the mock service
	handler := NewSurahHandler(mockService)

	// Create a new Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	c.Request = req

	// Call the method to be tested
	handler.GetListSurah(c)

	// Assert the results
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "success", response.Message)

	// Assert that the expected methods were called
	mockService.AssertExpectations(t)
}

func TestSurahHandler_GetListSurah_Error(t *testing.T) {
	// Create a new mock service
	mockService := new(MockSurahService)

	// Expect a call to GetListSurah and return an error
	mockService.On("GetListSurah", mock.Anything).Return(nil, errors.New("service error"))

	// Create a new surahHandler with the mock service
	handler := NewSurahHandler(mockService)

	// Create a new Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	c.Request = req

	// Call the method to be tested
	handler.GetListSurah(c)

	// Assert the results
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, response.Status)
	assert.Equal(t, "service error", response.Message)

	// Assert that the expected methods were called
	mockService.AssertExpectations(t)
}


func TestSurahHandler_GetDetailSurah_Success(t *testing.T) {
	// Create a new mock service
	mockService := new(MockSurahService)

	// Create a sample surah detail
	surahDetail := []dto.DetailSurahResp{
		{
			ID:    1,
			SurahID: 1,
			Ayah: 1,
			Page: 1,
			QuarterHizb: 1,
			Juz: 1,
			Manzil: 1,
			Arabic: "Arabic",
			Kitabah: "Kitabah",
			Latin: "Latin",
			ArabicWords: []string{"ArabicWords"},
			Translation: "Translation",
			UpdatedAt: time.Now(),
			Audio: "Audio",
			Surah: dto.SurahResp{
				ID: 1,
				Arabic: "Arabic",
				Latin: "Latin",
				Transliteration: "Transliteration",
				Translation: "Translation",
				NumAyah: 1,
				Page: 1,
				Location: "Location",
				UpdatedAt: time.Now(),
			},
		},
	}

	// Expect a call to GetSurahDetail and return the sample surah detail
	mockService.On("GetSurahDetail", mock.Anything, 1, 0, 10).Return(surahDetail, nil)

	// Create a new surahHandler with the mock service
	handler := NewSurahHandler(mockService)

	// Create a new Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	c.Request = req

	// Call the method to be tested
	handler.GetDetailSurah(c)

	// Assert the results
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "success", response.Message)

	// Assert that the expected methods were called
	mockService.AssertExpectations(t)
}