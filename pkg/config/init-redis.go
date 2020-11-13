// +build !no_config_init,!no_config_redis

package config

import (
	"github.com/boxgo/box/pkg/config/source/redis"
)

func init() {
	if !bootOK {
		return
	}

	for _, cfg := range sourceConfigs {
		if cfg.name != "redis" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[cfg.idx] = redis.NewSource(redis.WithConfig(cfg.data)...)
	}
}
