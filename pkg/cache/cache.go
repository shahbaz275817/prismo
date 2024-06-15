package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	redis.Cmdable
	Process(ctx context.Context, cmd redis.Cmder) error
	Close() error
}

type Options struct {
	Address                                             string
	PoolSize                                            int
	DialTimeout, ReadTimeout, WriteTimeout, IdleTimeout time.Duration
	UniqueKeyExpireTime                                 time.Duration
}

func NewClient(opts Options) (Client, error) {
	if opts.Address == "" {
		return nil, errors.New("No address to connect")
	}
	client := newRedisClient(opts)

	_, err := client.Ping(context.Background()).Result()
	return client, err
}

func newRedisClient(opts Options) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         opts.Address,
		PoolSize:     opts.PoolSize,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
	})
}
