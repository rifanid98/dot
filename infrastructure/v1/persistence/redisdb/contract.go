package redisdb

import (
	"context"
	"time"
)

//go:generate mockery --name Client --filename client.go --output ./mocks
type Client interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (*string, error)
	Del(ctx context.Context, key string) error
	HSet(ctx context.Context, key string, value map[string]interface{}) error
	Publish(ctx context.Context, channel string, data string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Ping(ctx context.Context) (string, error)
}
