package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/adexcell/shortener/internal/adapters/postgres"
	"github.com/wb-go/wbf/redis"
)

type ShortenerUsecase struct {
	postgres *postgres.URLPostgres
	redis    *redis.Client
}

func NewShortenerUsecase(p *postgres.URLPostgres, r *redis.Client) *ShortenerUsecase {
	return &ShortenerUsecase{postgres: p, redis: r}
}

// Shorten генерирует код и сохраняет в БД
func (u *ShortenerUsecase) Shorten(ctx context.Context, longURL string) (string, error) {
	b := make([]byte, 4)
	rand.Read(b)
	shortCode := base64.URLEncoding.EncodeToString(b)[:6]

	err := u.postgres.Save(ctx, shortCode, longURL)
	if err != nil {
		return "", err
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

	_ = u.redis.SetWithExpiration(ctx, shortCode, longURL, 24*time.Hour)
	
	go u.postgres.SaveClick(context.Background(), shortCode, ip, userAgent)

	return longURL, nil
}

func (u *ShortenerUsecase) GetStats(ctx context.Context, shortCode string) (int, error) {
	return u.postgres.GetAnalytics(ctx, shortCode)
}
