package redislocker

import "time"

var (
	Default = StdConfig("default").Build()
)

func Lock(key string, expire time.Duration) (bool, error) {
	return Default.Lock(key, expire)
}
func IsLocked(key string) (bool, error) {
	return Default.IsLocked(key)
}
func UnLock(key string) error {
	return Default.UnLock(key)
}
