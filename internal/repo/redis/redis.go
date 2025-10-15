package redisrepo

import (
	"context"
	"crypto/tls"
	"net/url"
	"strings"

	"github.com/redis/go-redis/v9"
)

func Open(addr string) *redis.Client {
	// Поддержка redis://[:password]@host:port и rediss:// (TLS)
	if strings.HasPrefix(addr, "redis://") || strings.HasPrefix(addr, "rediss://") {
		u, _ := url.Parse(addr)
		pass, _ := u.User.Password()
		opts := &redis.Options{
			Addr:     u.Host, // host:port
			Password: pass,
		}
		if u.Scheme == "rediss" {
			opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
		return redis.NewClient(opts)
	}
	// Обычный host:port без пароля
	return redis.NewClient(&redis.Options{Addr: addr})
}

func Ping(ctx context.Context, r *redis.Client) error { return r.Ping(ctx).Err() }

