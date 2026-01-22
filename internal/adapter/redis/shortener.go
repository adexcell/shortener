package redis

import (
	"context"
	"time"

	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/pkg/redis"
)

type Redis struct {
	redis *redis.RDB
}

func New(cfg redis.Config) domain.ShortenerRedis {
	redis := redis.New(cfg)
	return &Redis{redis: redis}
}

func (r *Redis) SetWithExpiration(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	return r.redis.SetWithExpiration(ctx, key, value, expiration)
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key)
}

func (r *Redis) Close() error {
	return r.redis.Close()
}
