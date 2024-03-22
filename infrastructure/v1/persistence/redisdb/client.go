package redisdb

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type clientImpl struct {
	client *redis.Client
}

func (redis *clientImpl) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return redis.client.Set(ctx, key, value, expiration).Err()
}

func (redis *clientImpl) Get(ctx context.Context, key string) (*string, error) {
	cmd := redis.client.Get(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	value := cmd.Val()
	return &value, nil
}

func (redis *clientImpl) Del(ctx context.Context, key string) error {
	return redis.client.Del(ctx, key).Err()
}

func (redis *clientImpl) HSet(ctx context.Context, key string, value map[string]interface{}) error {
	return redis.client.HMSet(ctx, key, value).Err()
}

func (redis *clientImpl) Publish(ctx context.Context, channel string, data string) error {
	return redis.client.Publish(ctx, channel, data).Err()
}

func (redis *clientImpl) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return redis.client.Expire(ctx, key, expiration).Err()
}

func (redis *clientImpl) Ping(ctx context.Context) (string, error) {
	return redis.client.Ping(ctx).Result()
}
