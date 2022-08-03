package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boxgo/box/v2/config"
)

type (
	Cache interface {
		Get(context.Context, string, interface{}) error
		Set(context.Context, string, interface{}, time.Duration) error
		Clear(context.Context, string) error
	}
)

var (
	ErrCacheMiss = errors.New("cache: key is missing")
)

func UnifiedKey(key string) string {
	return fmt.Sprintf("%s.cache.%s", config.ServiceName(), key)
}
