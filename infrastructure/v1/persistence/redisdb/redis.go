package redisdb

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"dot/config"
	"dot/pkg/util"
)

var log = util.NewLogger()

func New(config *config.RedisConfig) (*clientImpl, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Error(context.Background(), "failed to ping redis : %v", err)
		panic(err)
	}

	return &clientImpl{client}, nil
}

func NewRedisDbMayang(config *config.RedisConfig) (*clientImpl, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Error(context.Background(), "failed to ping redis : %v", err)
		panic(err)
	}

	return &clientImpl{client}, nil
}
