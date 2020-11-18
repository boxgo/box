// +build !no_config_init,!no_config_mongo

package config

import (
	"github.com/boxgo/box/pkg/config/source/mongodb"
)

func init() {
	if !bootOK {
		return
	}

	for _, cfg := range sourceConfigs {
		if cfg.name != "mongodb" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[cfg.idx] = mongodb.NewSource(mongodb.WithConfig(cfg.data)...)
	}
}
