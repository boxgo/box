// +build !no_config_init,!no_config_mongo

package config

import (
	"github.com/boxgo/box/pkg/config/source/mongodb"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "mongodb" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[idx] = mongodb.NewSource(mongodb.WithConfig(cfg.data)...)
	}
}
