// Package usecase содержит основную бизнес-логику.
package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/pkg/log"
	"github.com/adexcell/shortener/pkg/postgres"
)

type ShortenerUsecase struct {
	log      log.Log
	postgres domain.ShortenerPostgres
	redis    domain.ShortenerRedis
	ttl      time.Duration
	statsCh  chan domain.Stats
	wg       sync.WaitGroup
	mu       sync.RWMutex
	closed   bool
}

func New(p domain.ShortenerPostgres, r domain.ShortenerRedis, l log.Log, t time.Duration) domain.ShortenerUsecase {
	u := &ShortenerUsecase{
		log:      l,
		postgres: p,
		redis:    r,
		ttl:      t,
		statsCh:  make(chan domain.Stats, 1000),
	}

	u.wg.Add(1)
	go u.runAnalyticsWorker()

	return u
}

// Shorten генерирует код и сохраняет в БД
func (u *ShortenerUsecase) Shorten(ctx context.Context, shortCode, longURL string) (string, error) {
	if shortCode == "" {
		b := make([]byte, 4)
		rand.Read(b)
		shortCode = base64.URLEncoding.EncodeToString(b)[:6]
	}

	err := u.postgres.Save(ctx, shortCode, longURL)
	if err != nil {
		return "", postgres.PostgresErr(err)
	}

	if err := u.redis.SetWithExpiration(ctx, shortCode, longURL, u.ttl); err != nil {
		u.log.Error().Err(err).Str("code", shortCode).Msg("failed to save click analytics in redis")
	}

	return shortCode, nil
}

// GetOriginal ищет полную ссылку по коду
func (u *ShortenerUsecase) GetOriginal(ctx context.Context, shortCode, ip, userAgent string) (string, error) {
	longURL, err := u.redis.Get(ctx, shortCode)
	if err != nil {
		longURL, err = u.postgres.GetLongURL(ctx, shortCode)
		if err != nil {
			return "", fmt.Errorf("failed to get long url from db: %w", err)
		}

		if err := u.redis.SetWithExpiration(ctx, shortCode, longURL, u.ttl); err != nil {
			u.log.Error().Err(err).Str("code", shortCode).Msg("failed to save click analytics in redis")
		}
	}

	stats := domain.Stats{
		ShortCode: shortCode,
		IP:        ip,
		UserAgent: userAgent,
	}

	u.mu.RLock()
	defer u.mu.RUnlock()
	if u.closed {
		return longURL, nil
	}
	select {
	case u.statsCh <- stats:
		// успешно отправили
	default:
		u.log.Warn().Str("code", shortCode).Msg("analytics channel full, dropping stat")
	}

	return longURL, nil
}

func (u *ShortenerUsecase) GetStats(ctx context.Context, shortCode string) (domain.Stats, error) {
	return u.postgres.GetDetailedStats(ctx, shortCode)
}

func (u *ShortenerUsecase) runAnalyticsWorker() {
	defer u.wg.Done()

	u.log.Info().Msg("analytics worker started")

	for stats := range u.statsCh {
		err := u.postgres.SaveClick(context.Background(), stats.ShortCode, stats.IP, stats.UserAgent)
		if err != nil {
			u.log.Error().Err(err).Str("code", stats.ShortCode).Msg("failed to save click analytics in postgres")
		}
	}
}

func (u *ShortenerUsecase) Close() error {
	u.mu.Lock()
	u.closed = true
	u.mu.Unlock()

	close(u.statsCh)

	u.wg.Wait()
	return nil
}
