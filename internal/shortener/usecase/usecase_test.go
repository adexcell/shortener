package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/adexcell/shortener/internal/shortener/dto"
	"github.com/adexcell/shortener/internal/shortener/worker"
)

func TestShortenUsecase_CreateShorten(t *testing.T) {
	ctx := context.Background()
	mockPostgres := new(MockPostgres)
	mockRedis := new(MockRedis)
	mockRedisWriter := worker.NewAsyncRedisWriter(mockRedis)
	
	u := &ShortenUsecase{
		postgres:         mockPostgres,
		redis:            mockRedis,
		asyncRedisWriter: mockRedisWriter,
	}

	t.Run("Success with custom alias", func(t *testing.T) {
		input := dto.CreateShortenInput{
			OriginalURL: "https://google.com",
			ShortenCode: "myalias",
		}

		mockPostgres.On("SaveURL", ctx, mock.MatchedBy(func(s interface{}) bool {
			shorten := s.(domain.Shortener)
			return shorten.ShortCode == "myalias" && shorten.OriginalURL == "https://google.com"
		})).Return(nil).Once()

		output, err := u.CreateShorten(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, "myalias", output.ShortenCode)
		mockPostgres.AssertExpectations(t)
	})

	t.Run("Success with generated code", func(t *testing.T) {
		input := dto.CreateShortenInput{
			OriginalURL: "https://google.com",
		}

		mockPostgres.On("SaveURL", ctx, mock.Anything).Return(nil).Once()

		output, err := u.CreateShorten(ctx, input)

		assert.NoError(t, err)
		assert.NotEmpty(t, output.ShortenCode)
		assert.Len(t, output.ShortenCode, 6)
		mockPostgres.AssertExpectations(t)
	})

	t.Run("Postgres Error", func(t *testing.T) {
		input := dto.CreateShortenInput{
			OriginalURL: "https://google.com",
			ShortenCode: "error",
		}

		mockPostgres.On("SaveURL", ctx, mock.Anything).Return(assert.AnError).Once()

		_, err := u.CreateShorten(ctx, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "postgres.SaveURL")
		mockPostgres.AssertExpectations(t)
	})
}

func TestShortenUsecase_GetOriginalURL(t *testing.T) {
	ctx := context.Background()
	mockPostgres := new(MockPostgres)
	mockRedis := new(MockRedis)
	mockRedisWriter := worker.NewAsyncRedisWriter(mockRedis)
	
	u := &ShortenUsecase{
		postgres:         mockPostgres,
		redis:            mockRedis,
		asyncRedisWriter: mockRedisWriter,
		statsCh:          make(chan dto.RedirectInput, 10),
	}

	t.Run("Cache Hit", func(t *testing.T) {
		input := dto.RedirectInput{ShortenCode: "hit"}
		url := "https://google.com"
		mockRedis.On("Get", ctx, "hit").Return(&url, nil).Once()

		res, err := u.GetOriginalURL(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, url, res)
		mockRedis.AssertExpectations(t)
	})

	t.Run("Cache Miss, DB Hit", func(t *testing.T) {
		input := dto.RedirectInput{ShortenCode: "miss"}
		url := "https://google.com"
		mockRedis.On("Get", ctx, "miss").Return(nil, assert.AnError).Once()
		mockPostgres.On("GetOriginalURL", ctx, "miss").Return(url, nil).Once()

		res, err := u.GetOriginalURL(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, url, res)
		mockPostgres.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
	})
}

func TestShortenUsecase_GetAnalytics(t *testing.T) {
	ctx := context.Background()
	mockPostgres := new(MockPostgres)
	
	u := &ShortenUsecase{
		postgres: mockPostgres,
	}

	t.Run("Success", func(t *testing.T) {
		input := dto.GetAnalyticsInput{ShortCode: "code"}
		expectedStats := domain.Stats{
			ShortCode:   "code",
			TotalClicks: 10,
			ByBrowser:   map[string]int{"Chrome": 10},
		}
		mockPostgres.On("GetDetailedStats", ctx, mock.Anything).Return(expectedStats, nil).Once()

		res, err := u.GetAnalytics(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, 10, res.TotalClicks)
		assert.Equal(t, 10, res.ByBrowser["Chrome"])
		mockPostgres.AssertExpectations(t)
	})
}

