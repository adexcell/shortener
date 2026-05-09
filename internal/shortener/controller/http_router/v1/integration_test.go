package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/shortener/internal/shortener/dto"
)

type MockUsecase struct {
	mock.Mock
}

func (m *MockUsecase) CreateShorten(ctx context.Context, input dto.CreateShortenInput) (dto.CreateShortenOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(dto.CreateShortenOutput), args.Error(1)
}

func (m *MockUsecase) GetOriginalURL(ctx context.Context, input dto.RedirectInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}

func (m *MockUsecase) GetAnalytics(ctx context.Context, input dto.GetAnalyticsInput) (dto.GetAnalyticsOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(dto.GetAnalyticsOutput), args.Error(1)
}

func TestHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUc := new(MockUsecase)
	handler := New(mockUc)
	
	r := ginext.New("test")
	api := r.Group("/api/v1")
	{
		api.POST("/shorten", handler.CreateShorten)
		api.GET("/s/:short_url", handler.RedirectShortLink)
		api.GET("/analytics/:short_url", handler.GetAnalytics)
	}

	t.Run("POST /shorten - Success", func(t *testing.T) {
		input := dto.CreateShortenInput{
			OriginalURL: "https://google.com",
			ShortenCode: "myalias",
		}
		mockUc.On("CreateShorten", mock.Anything, mock.Anything).Return(dto.CreateShortenOutput{ShortenCode: "myalias"}, nil).Once()

		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CreateShortenOutput
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "myalias", resp.ShortenCode)
	})

	t.Run("GET /s/:short_url - Success", func(t *testing.T) {
		mockUc.On("GetOriginalURL", mock.Anything, mock.MatchedBy(func(in dto.RedirectInput) bool {
			return in.ShortenCode == "qwerty"
		})).Return("https://google.com", nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/s/qwerty", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "https://google.com", w.Header().Get("Location"))
	})

	t.Run("GET /analytics/:short_url - Success", func(t *testing.T) {
		mockUc.On("GetAnalytics", mock.Anything, mock.Anything).Return(dto.GetAnalyticsOutput{TotalClicks: 5}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/analytics/qwerty", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.GetAnalyticsOutput
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 5, resp.TotalClicks)
	})
}
