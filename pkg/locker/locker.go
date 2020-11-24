package locker

import (
	"context"
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/config"
)

type (
	MutexLocker interface {
		Lock(context.Context, string, time.Duration) (bool, error)
		IsLocked(context.Context, string) (bool, error)
		UnLock(context.Context, string) error
	}
)

func UnifiedKey(key string) string {
	return fmt.Sprintf("%s.locker.%s", config.ServiceName(), key)
}
