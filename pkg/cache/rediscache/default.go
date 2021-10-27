package rediscache

import (
	"context"
	"time"
)

var (
	Default = StdConfig("default").Build()
)

func Get(ctx context.Context, key string, val interface{}) error {
	return Default.Get(ctx, key, val)
}

func Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	return Default.Set(ctx, key, val, ttl)
}
