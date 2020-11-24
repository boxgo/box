package locker

import (
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/config"
)

type (
	MutexLocker interface {
		Lock(string, time.Duration) (bool, error)
		IsLocked(string) (bool, error)
		UnLock(string) error
	}
)

func UnifiedKey(key string) string {
	return fmt.Sprintf("%s.locker.%s", config.ServiceName(), key)
}
