// +build !no_config_init,!no_config_file

package config

import (
	file2 "github.com/boxgo/box/v2/config/source/file"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "file" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[idx] = file2.NewSource(file2.WithConfig(cfg.data)...)
	}
}
