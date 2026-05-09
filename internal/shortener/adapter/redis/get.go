package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Get retrieves the popular shorten URLs from Redis by its idempotency key.
func (r *Redis) Get(ctx context.Context, idempotencyKey string) (*string, error) {
	key := idempotencyPrefix + idempotencyKey

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, err
		}
	log.Error().Err(err).Msg("redis.Get failed")
	return nil, err
	}

	return &val, nil
}
