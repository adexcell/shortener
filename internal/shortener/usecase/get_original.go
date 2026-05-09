package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adexcell/shortener/internal/shortener/dto"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (u *ShortenUsecase) GetOriginalURL(ctx context.Context, input dto.RedirectInput) (string, error) {
	var originalURL string

	// get from cache
	originalURLPntr, err := u.redis.Get(ctx, input.ShortenCode)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Debug().Err(err).Msg("cache miss")
		} else {
			log.Error().Err(err).Msg("redis unavailable")
		}

		// get from db
		originalURL, err = u.postgres.GetOriginalURL(ctx, input.ShortenCode)
		if err != nil {
			return "", fmt.Errorf("postgres.GetLongURL: %w", err)
		}
	}
	if err == nil && originalURLPntr != nil {
		log.Debug().Msg("successful cache hit")
		originalURL = *originalURLPntr
	}

	// save to cache
	task := dto.ShortenTask{
		Shorten:     input.ShortenCode,
		OriginalURL: originalURL,
	}
	u.asyncRedisWriter.Send(ctx, task)

	// save stats
	if !u.isClosed.Load() {
		select {
		case u.statsCh <- input:
		default:
			log.Warn().Msg("usecase stats channel full, dropping stat")
		}
	}

	return originalURL, nil
}

func (u *ShortenUsecase) runAnalyticsWorker() {
	defer u.wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Info().Msg("usecase worker started")
	for stats := range u.statsCh {
		u.SaveClick(ctx, stats)
	}

}

func (u *ShortenUsecase) Stop() {
	log.Debug().Bool("usecase.isClosed", u.isClosed.Load())
	if u.isClosed.CompareAndSwap(false, true) {
		close(u.statsCh)
	}
	u.wg.Wait()
}
