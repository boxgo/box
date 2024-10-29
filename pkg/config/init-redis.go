//go:build !no_config_init && !no_config_redis
// +build !no_config_init,!no_config_redis

package config

import (
	"github.com/boxgo/box/pkg/config/source/redis"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "redis" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[idx] = redis.NewSource(
			append(
				redis.WithConfig(cfg.data),
				redis.WithPrefix(bootCfg.Service),
			)...,
		)
	}
}
