package redisrepo

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func Open(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr})
}

func Ping(ctx context.Context, r *redis.Client) error { return r.Ping(ctx).Err() }
