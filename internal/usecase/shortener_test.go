package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/internal/usecase"
	"github.com/adexcell/shortener/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

const (
	TTL = 24 * time.Hour
)

type MockPostgres struct {
	mock.Mock
}

func (m *MockPostgres) Save(ctx context.Context, shortCode, longURL string) error {
	args := m.Called(ctx, shortCode, longURL)
	return args.Error(0)
}

func (m *MockPostgres) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Error(1)
}

func (m *MockPostgres) SaveClick(ctx context.Context, shortCode, ip, userAgent string) error {
	args := m.Called(ctx, shortCode, ip, userAgent)
	return args.Error(0)
}

func (m *MockPostgres) GetDetailedStats(ctx context.Context, shortCode string) (domain.Stats, error) {
	args := m.Called(ctx, shortCode)
	return args.Get(0).(domain.Stats), args.Error(1)
}

func (m *MockPostgres) Close() error {
	return nil
}

type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedis) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) Close() error {
	return nil
}

// --- Tests ---

func TestShortenerUsecase_Shorten(t *testing.T) {
	ctx := context.Background()
	mockPg := new(MockPostgres)
	mockRedis := new(MockRedis)
	log := log.New()

	uc := usecase.New(mockPg, mockRedis, log, TTL)
	longURL := "https://example.com"

	t.Run("success", func(t *testing.T) {
		// Expectation: Save to postgres, then save to redis
		mockPg.On("Save", ctx, mock.AnythingOfType("string"), longURL).Return(nil).Once()
		mockRedis.On("SetWithExpiration", ctx, mock.AnythingOfType("string"), longURL, 24*time.Hour).Return(nil).Once()

		code, err := uc.Shorten(ctx, "", longURL)

		assert.NoError(t, err)
		assert.NotEmpty(t, code)
		assert.Len(t, code, 6)
		mockPg.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
	})

	t.Run("postgres error", func(t *testing.T) {
		mockPg.On("Save", ctx, mock.AnythingOfType("string"), longURL).Return(errors.New("db error")).Once()

		code, err := uc.Shorten(ctx, "", longURL)

		assert.Error(t, err)
		assert.Empty(t, code)
		mockPg.AssertExpectations(t)
	})
}

func TestShortenerUsecase_GetOriginal(t *testing.T) {
	mockPg := new(MockPostgres)
	mockRedis := new(MockRedis)
	log := log.New()
	ctx := context.Background()

	uc := usecase.New(mockPg, mockRedis, log, TTL)
	shortCode := "abcdef"
	longURL := "https://example.com"
	ip := "127.0.0.1"
	ua := "Go-Test"

	t.Run("redis hit", func(t *testing.T) {
		mockRedis.On("Get", ctx, shortCode).Return(longURL, nil).Once()
		// SaveClick is called in a goroutine, so we might not be able to assert it reliably without sync.
		// However, the mock might catch it if we add a wait or sleep, but for unit test speed we often ignore async calls or just mock them loosely.
		// For this test, we accept if it's called or not, BUT mocks are strict.
		// If the code calls mock, we MUST expect it or use .Maybe()
		// Since it's async context.Background(), the mocking might happen after test finishes.
		// We'll define it as .Maybe() or use a separate mock instance if we wanted strict async testing.
		// Let's use .Maybe() and .Return(nil) to prevent panic if it is called.
		mockPg.On("SaveClick", mock.Anything, shortCode, ip, ua).Return(nil).Maybe()

		url, err := uc.GetOriginal(ctx, shortCode, ip, ua)

		assert.NoError(t, err)
		assert.Equal(t, longURL, url)
		mockRedis.AssertExpectations(t)
	})

	t.Run("redis miss, postgres hit", func(t *testing.T) {
		mockRedis.On("Get", ctx, shortCode).Return("", errors.New("not found")).Once()
		mockPg.On("GetLongURL", ctx, shortCode).Return(longURL, nil).Once()
		mockRedis.On("SetWithExpiration", ctx, shortCode, longURL, 24*time.Hour).Return(nil).Once()
		mockPg.On("SaveClick", mock.Anything, shortCode, ip, ua).Return(nil).Maybe()

		url, err := uc.GetOriginal(ctx, shortCode, ip, ua)

		assert.NoError(t, err)
		assert.Equal(t, longURL, url)
		mockPg.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
	})

	t.Run("not found everywhere", func(t *testing.T) {
		mockRedis.On("Get", ctx, shortCode).Return("", errors.New("not found")).Once()
		mockPg.On("GetLongURL", ctx, shortCode).Return("", errors.New("not found")).Once()

		url, err := uc.GetOriginal(ctx, shortCode, ip, ua)

		assert.Error(t, err) // The usecase returns an error when URL is not found
		assert.Empty(t, url)
		mockPg.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
	})
}
