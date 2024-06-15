package cache

import (
	"github.com/shahbaz275817/prismo/pkg/config"
)

func NewCacheConfig() Options {
	return Options{
		Address:             config.MustGetString("REDIS_CACHE_HOST"),
		PoolSize:            config.MustGetInt("REDIS_CACHE_MAX_POOL_SIZE"),
		DialTimeout:         config.MustGetTimeoutInMS("REDIS_CACHE_CONNECT_TIMEOUT_MS"),
		ReadTimeout:         config.MustGetTimeoutInMS("REDIS_CACHE_QUERY_TIMEOUT_MS"),
		WriteTimeout:        config.MustGetTimeoutInMS("REDIS_CACHE_QUERY_TIMEOUT_MS"),
		IdleTimeout:         config.MustGetTimeoutInMS("REDIS_CACHE_MAX_IDLE_TIMEOUT_MS"),
		UniqueKeyExpireTime: config.MustGetTimeoutInMS("REDIS_CACHE_UNIQUE_KEY_EXPIRE_TIME_MS"),
	}
}
