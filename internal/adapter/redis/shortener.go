package redis

import (
	"context"
	"time"

	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/pkg/redis"
)

type ShortenerRedis struct {
	redis *redis.RDB
}

func NewShortenerRedis(cfg redis.Config) domain.ShortenerRedis {
	redis := redis.NewRedis(cfg)
	return &ShortenerRedis{redis: redis}
}

func (r *ShortenerRedis) SetWithExpiration(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	return r.redis.SetWithExpiration(ctx, key, value, expiration)
}

func (r *ShortenerRedis) Get(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key)
}
