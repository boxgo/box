package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/config"
)

type (
	Cache interface {
		Get(context.Context, string, interface{}) error
		Set(context.Context, string, interface{}, time.Duration) error
	}
)

func UnifiedKey(key string) string {
	return fmt.Sprintf("%s.cache.%s", config.ServiceName(), key)
}
