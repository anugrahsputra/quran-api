package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anugrahsputra/go-quran-api/domain/model"
	"github.com/blevesearch/bleve/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuranSearchRepository struct {
	mock.Mock
}

func (m *MockQuranSearchRepository) Index(ayahs []model.Ayah) error {
	args := m.Called(ayahs)
	return args.Error(0)
}

func (m *MockQuranSearchRepository) Search(query string, page, limit int) (*bleve.SearchResult, error) {
	args := m.Called(query, page, limit)
	return args.Get(0).(*bleve.SearchResult), args.Error(1)
}

func (m *MockQuranSearchRepository) GetDocument(id string) (map[string]any, error) {
	args := m.Called(id)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockQuranSearchRepository) GetDocCount() (uint64, error) {
	args := m.Called()
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockQuranSearchRepository) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockQuranSearchRepository)
	h := NewHealthHandler(mockRepo)

	r := gin.Default()
	r.GET("/ping", h.Ping)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"pong"}`, w.Body.String())
}
