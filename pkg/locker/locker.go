package locker

import (
	"context"
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/config"
)

type (
	MutexLocker interface {
		Lock(context.Context, string, time.Duration) (success bool, err error)
		IsLocked(context.Context, string) (locked bool, err error)
		UnLock(context.Context, string) error
	}
)

func UnifiedKey(key string) string {
	return fmt.Sprintf("%s.locker.%s", config.ServiceName(), key)
}
