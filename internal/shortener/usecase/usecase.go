package usecase

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/adexcell/shortener/internal/shortener/dto"
	"github.com/adexcell/shortener/internal/shortener/worker"
)

type Postgres interface {
	SaveURL(ctx context.Context, shorten domain.Shortener) error
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	SaveClick(ctx context.Context, stats domain.Stats) error
	GetDetailedStats(ctx context.Context, stats domain.Stats) (domain.Stats, error)
}

type Redis interface {
	Get(ctx context.Context, idempotencyKey string) (*string, error)
}

type ShortenUsecase struct {
	postgres         Postgres
	redis            Redis
	asyncRedisWriter *worker.AsyncRedisWriter
	statsCh          chan dto.RedirectInput
	isClosed         atomic.Bool
	wg               sync.WaitGroup
}

func New(ctx context.Context, postgres Postgres, redis Redis, asyncRedisWriter *worker.AsyncRedisWriter) *ShortenUsecase {
	usecase := &ShortenUsecase{
		postgres:         postgres,
		redis:            redis,
		asyncRedisWriter: asyncRedisWriter,
		statsCh:          make(chan dto.RedirectInput, 1000),
	}

	usecase.wg.Add(1)
	go usecase.runAnalyticsWorker()

	return usecase
}
