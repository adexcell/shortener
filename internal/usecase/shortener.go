// Package usecase содержит основную бизнес-логику.
package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/postgres"
)

type ShortenerUsecase struct {
	postgres domain.ShortenerPostgres
	redis    domain.ShortenerRedis
	log      logger.Log
}

func NewShortenerUsecase(
	p domain.ShortenerPostgres, 
	r domain.ShortenerRedis,
	l logger.Log,
	) domain.ShortenerUsecase {
	return &ShortenerUsecase{
		postgres: p, 
		redis: r,
		log: l,
	}
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

	_ = u.redis.SetWithExpiration(ctx, shortCode, longURL, 24*time.Hour)
	return shortCode, nil
}

// GetOriginal ищет полную ссылку по коду
func (u *ShortenerUsecase) GetOriginal(ctx context.Context, shortCode, ip, userAgent string) (string, error) {
	longURL, err := u.redis.Get(ctx, shortCode)
	if err == nil && longURL != "" {
		go u.postgres.SaveClick(context.Background(), shortCode, ip, userAgent)
		return longURL, nil
	}

	longURL, err = u.postgres.GetLongURL(ctx, shortCode)
	if err != nil {
		return "", nil
	}

	err = u.redis.SetWithExpiration(ctx, shortCode, longURL, 24*time.Hour)
	if err != nil {
			u.log.Error().Err(err).Str("code", shortCode).Msg("failed to save click analytics in redis")
		}

	go func() {
		err := u.postgres.SaveClick(context.Background(), shortCode, ip, userAgent)
		if err != nil {
			u.log.Error().Err(err).Str("code", shortCode).Msg("failed to save click analytics in postgres")
		}
	}()

	return longURL, nil
}

func (u *ShortenerUsecase) GetStats(ctx context.Context, shortCode string) (domain.Stats, error) {
	return u.postgres.GetDetailedStats(ctx, shortCode)
}
