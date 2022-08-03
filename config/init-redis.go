// +build !no_config_init,!no_config_redis

package config

import (
	redis2 "github.com/boxgo/box/v2/config/source/redis"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "redis" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[idx] = redis2.NewSource(
			append(
				redis2.WithConfig(cfg.data),
				redis2.WithPrefix(bootCfg.Name),
			)...,
		)
	}
}
