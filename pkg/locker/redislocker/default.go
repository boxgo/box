package redislocker

import (
	"context"
	"time"
)

var (
	Default = StdConfig("default").Build()
)

func Lock(ctx context.Context, key string, expire time.Duration) (bool, error) {
	return Default.Lock(ctx, key, expire)
}
func IsLocked(ctx context.Context, key string) (bool, error) {
	return Default.IsLocked(ctx, key)
}
func UnLock(ctx context.Context, key string) error {
	return Default.UnLock(ctx, key)
}
