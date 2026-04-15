package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(redisURL string) (*Redis, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("store.NewRedis: %w", err)
	}
	client := redis.NewClient(opts)
	return &Redis{client: client}, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// IncrRateLimit increments the rate limit counter for a client and returns the current count.
// The key expires after 1 minute.
func (r *Redis) IncrRateLimit(ctx context.Context, clientID string) (int64, error) {
	bucket := time.Now().UTC().Format("2006-01-02T15:04")
	key := fmt.Sprintf("verify:rl:%s:%s", clientID, bucket)

	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("store.IncrRateLimit: %w", err)
	}

	if count == 1 {
		r.client.Expire(ctx, key, 2*time.Minute)
	}

	return count, nil
}
