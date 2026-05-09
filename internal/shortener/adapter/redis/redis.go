package redis

import (
	"time"

	"github.com/adexcell/shortener/pkg/redis"
)

var (
	idempotencyPrefix = "shorten:idempotency:"
	ttl               = time.Hour
)

// Redis implements the usecase.Redis interface using Redis.
type Redis struct {
	client *redis.Client
}

// New creates a new instance of the Redis adapter.
func New(client *redis.Client) *Redis {
	return &Redis{client: client}
}
