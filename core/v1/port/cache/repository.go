package cache

import (
	"dot/core"
	"time"
)

//go:generate mockery --name CacheRepository --filename cache_repository.go --output ./mocks
type CacheRepository interface {
	Set(ic *core.InternalContext, key string, value string, expiration *time.Duration) *core.CustomError
	HSet(ic *core.InternalContext, key string, value map[string]interface{}, expiration time.Duration) *core.CustomError
	Get(ic *core.InternalContext, key string) (string, *core.CustomError)
	Delete(ic *core.InternalContext, key string) *core.CustomError
	Publish(ic *core.InternalContext, channel string, data string) *core.CustomError
}
