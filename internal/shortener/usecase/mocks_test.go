package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/adexcell/shortener/internal/shortener/domain"
)

type MockPostgres struct {
	mock.Mock
}

func (m *MockPostgres) SaveURL(ctx context.Context, shorten domain.Shortener) error {
	args := m.Called(ctx, shorten)
	return args.Error(0)
}

func (m *MockPostgres) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Error(1)
}

func (m *MockPostgres) SaveClick(ctx context.Context, stats domain.Stats) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockPostgres) GetDetailedStats(ctx context.Context, stats domain.Stats) (domain.Stats, error) {
	args := m.Called(ctx, stats)
	return args.Get(0).(domain.Stats), args.Error(1)
}

type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Get(ctx context.Context, idempotencyKey string) (*string, error) {
	args := m.Called(ctx, idempotencyKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockRedis) Set(ctx context.Context, key, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

