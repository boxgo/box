// +build !no_config_init,!no_config_mongo

package config

import (
	mongodb2 "github.com/boxgo/box/v2/config/source/mongodb"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "mongodb" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[idx] = mongodb2.NewSource(
			append(
				mongodb2.WithConfig(cfg.data),
				mongodb2.WithService(bootCfg.Name),
			)...,
		)
	}
}
