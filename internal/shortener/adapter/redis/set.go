package redis

import (
	"context"
	"fmt"
)

// Set stores the popular shorten URLs in Redis with an idempotency key.
func (r *Redis) Set(ctx context.Context, idempotencyKey, value string) error {
	key := idempotencyPrefix + idempotencyKey

	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}
