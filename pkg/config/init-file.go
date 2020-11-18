// +build !no_config_init,!no_config_file

package config

import (
	"github.com/boxgo/box/pkg/config/source/file"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "file" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[idx] = file.NewSource(file.WithConfig(cfg.data)...)
	}
}
