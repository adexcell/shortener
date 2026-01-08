package controller_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adexcell/shortener/internal/controller"
	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockUsecase struct {
	mock.Mock
}

func (m *MockUsecase) Shorten(ctx context.Context, shortCode, longURL string) (string, error) {
	args := m.Called(ctx, shortCode, longURL)
	return args.String(0), args.Error(1)
}

func (m *MockUsecase) GetOriginal(ctx context.Context, shortCode, ip, userAgent string) (string, error) {
	args := m.Called(ctx, shortCode, ip, userAgent)
	return args.String(0), args.Error(1)
}

func (m *MockUsecase) GetStats(ctx context.Context, shortCode string) (domain.Stats, error) {
	args := m.Called(ctx, shortCode)
	return args.Get(0).(domain.Stats), args.Error(1)
}

func (m *MockUsecase) Close() error {
	// TODO:
	return nil
}

// --- Setup ---

func setupRouter() *router.Router {
	// Use test mode to suppress debug logs
	return router.NewRouter(router.Config{GinMode: "test"})
}

// --- Tests ---

func TestHandler_PostShortURL(t *testing.T) {
	log := logger.NewLogger()

	t.Run("success", func(t *testing.T) {
		mockUC := new(MockUsecase)
		r := setupRouter()
		h := controller.NewShortenHandler(mockUC, log)
		h.Register(r)

		inputBody := `{"url": "https://example.com"}`
		expectedCode := "abcdef"

		// Expectation
		mockUC.On("Shorten", mock.Anything, "", "https://example.com").Return(expectedCode, nil)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(inputBody))
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), expectedCode)
		mockUC.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockUC := new(MockUsecase)
		r := setupRouter()
		h := controller.NewShortenHandler(mockUC, log)
		h.Register(r)

		// Request (empty body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(""))
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertNotCalled(t, "Shorten")
	})

	t.Run("internal error", func(t *testing.T) {
		mockUC := new(MockUsecase)
		r := setupRouter()
		h := controller.NewShortenHandler(mockUC, log)
		h.Register(r)

		inputBody := `{"url": "https://example.com", "alias": "custom"}`

		// Expectation
		mockUC.On("Shorten", mock.Anything, "custom", "https://example.com").Return("", errors.New("db fail"))

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(inputBody))
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestHandler_ConversionURL(t *testing.T) {
	log := logger.NewLogger()

	t.Run("redirect success", func(t *testing.T) {
		mockUC := new(MockUsecase)
		r := setupRouter()
		h := controller.NewShortenHandler(mockUC, log)
		h.Register(r)

		shortCode := "abc1234"
		longURL := "https://google.com"

		// Expectation
		mockUC.On("GetOriginal", mock.Anything, shortCode, mock.Anything, mock.Anything).Return(longURL, nil)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/s/"+shortCode, nil)
		// Mock ClientIP for validation (Gin uses RemoteAddr or headers)
		req.RemoteAddr = "127.0.0.1:12345"

		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusFound, w.Code) // 302
		assert.Equal(t, longURL, w.Header().Get("Location"))
		mockUC.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockUC := new(MockUsecase)
		r := setupRouter()
		h := controller.NewShortenHandler(mockUC, log)
		h.Register(r)

		shortCode := "missing"

		// Expectation
		mockUC.On("GetOriginal", mock.Anything, shortCode, mock.Anything, mock.Anything).Return("", errors.New("not found"))

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/s/"+shortCode, nil)
		req.RemoteAddr = "127.0.0.1:12345"

		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUC.AssertExpectations(t)
	})
}
