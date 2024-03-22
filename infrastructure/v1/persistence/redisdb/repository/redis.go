package repository

import (
	"dot/pkg/util"
	"time"

	"dot/core"
	"dot/infrastructure/v1/persistence/redisdb"
)

var log = util.NewLogger()

type cacheRepositoryImpl struct {
	client     redisdb.Client
	expiration time.Duration
}

func NewCacheRepository(client redisdb.Client) *cacheRepositoryImpl {
	return &cacheRepositoryImpl{
		client: client,
	}
}

func (redis *cacheRepositoryImpl) Set(ic *core.InternalContext, key string, value string, expiration *time.Duration) *core.CustomError {
	if expiration == nil {
		expiration = &redis.expiration
	}

	err := redis.client.Set(ic.ToContext(), key, value, *expiration)
	if err != nil {
		log.Error(ic.ToContext(), "failed to Set", err)
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func (redis *cacheRepositoryImpl) HSet(ic *core.InternalContext, key string, value map[string]interface{}, expiration time.Duration) *core.CustomError {
	err := redis.client.HSet(ic.ToContext(), key, value)
	if err != nil {
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	err = redis.client.Expire(ic.ToContext(), key, expiration)
	if err != nil {
		log.Error(ic.ToContext(), "failed to HSet", err)
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func (redis *cacheRepositoryImpl) Get(ic *core.InternalContext, key string) (string, *core.CustomError) {
	get, err := redis.client.Get(ic.ToContext(), key)
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", nil
		}

		log.Error(ic.ToContext(), "failed to Get", err)
		return "", &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}
	if get == nil {
		return "", nil
	}
	if *get == "0" {
		return "", nil
	}
	if *get == "null" {
		return "", nil
	}
	return *get, nil
}

func (redis *cacheRepositoryImpl) Delete(ic *core.InternalContext, key string) *core.CustomError {
	err := redis.client.Del(ic.ToContext(), key)
	if err != nil {
		log.Error(ic.ToContext(), "failed to Delete", err)
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func (redis *cacheRepositoryImpl) Publish(ic *core.InternalContext, channel string, data string) *core.CustomError {
	err := redis.client.Publish(ic.ToContext(), channel, data)
	if err != nil {
		log.Error(ic.ToContext(), "failed to Publish", err)
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}
