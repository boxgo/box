package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boxgo/box/v2/build"
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

// UnifiedKey namespace+name+key
func UnifiedKey(key string) string {
	return fmt.Sprintf("%s.%s.cache.%s", build.Namespace, build.Name, key)
}

// UnifiedKeyVer namespace+name+version+key
func UnifiedKeyVer(key string) string {
	return fmt.Sprintf("%s.%s.%s.cache.%s", build.Namespace, build.Name, build.Version, key)
}

// UnifiedKeyId namespace+name+version+id+key
func UnifiedKeyId(key string) string {
	return fmt.Sprintf("%s.%s.%s.%s.cache.%s", build.Namespace, build.Name, build.Version, build.ID, key)
}
