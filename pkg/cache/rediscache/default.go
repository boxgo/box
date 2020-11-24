package rediscache

import (
	"context"
	"time"
)

var (
	def = StdConfig("default").Build()
)

func Get(ctx context.Context, key string, val interface{}) error {
	return def.Get(ctx, key, val)
}

func Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	return def.Set(ctx, key, val, ttl)
}
