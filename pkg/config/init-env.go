// +build !no_config_init,!no_config_env

package config

import (
	"strings"

	"github.com/boxgo/box/pkg/config/source/env"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "env" || len(cfg.data) == 0 {
			continue
		}

		prefix := strings.ToUpper(bootCfg.Name)
		defaultSources[idx] = env.NewSource(env.WithStrippedPrefix(prefix))
	}
}
