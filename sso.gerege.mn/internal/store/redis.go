package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

type Redis struct {
	client valkey.Client
}

func NewRedis(redisURL string) (*Redis, error) {
	opts, err := valkey.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("store.NewRedis parse: %w", err)
	}
	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("store.NewRedis: %w", err)
	}
	return &Redis{client: client}, nil
}

func (r *Redis) Close() {
	r.client.Close()
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis.Set marshal: %w", err)
	}
	cmd := r.client.B().Set().Key(key).Value(string(data)).Ex(ttl).Build()
	return r.client.Do(ctx, cmd).Error()
}

func (r *Redis) Get(ctx context.Context, key string, dest any) error {
	cmd := r.client.B().Get().Key(key).Build()
	result, err := r.client.Do(ctx, cmd).ToString()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), dest)
}

func (r *Redis) Del(ctx context.Context, key string) error {
	cmd := r.client.B().Del().Key(key).Build()
	return r.client.Do(ctx, cmd).Error()
}

func (r *Redis) GetAndDel(ctx context.Context, key string, dest any) error {
	cmd := r.client.B().Getdel().Key(key).Build()
	result, err := r.client.Do(ctx, cmd).ToString()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), dest)
}

func (r *Redis) Incr(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	cmd := r.client.B().Incr().Key(key).Build()
	val, err := r.client.Do(ctx, cmd).AsInt64()
	if err != nil {
		return 0, err
	}
	if val == 1 {
		expCmd := r.client.B().Expire().Key(key).Seconds(int64(ttl.Seconds())).Build()
		r.client.Do(ctx, expCmd)
	}
	return val, nil
}

func (r *Redis) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	cmd := r.client.B().Set().Key(key).Value(value).Ex(ttl).Build()
	return r.client.Do(ctx, cmd).Error()
}

func (r *Redis) GetString(ctx context.Context, key string) (string, error) {
	cmd := r.client.B().Get().Key(key).Build()
	return r.client.Do(ctx, cmd).ToString()
}
