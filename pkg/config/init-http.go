//go:build !no_config_init && !no_config_http
// +build !no_config_init,!no_config_http

package config

import (
	"github.com/boxgo/box/pkg/config/source/http"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "http" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[idx] = http.NewSource(
			append(
				http.WithConfig(cfg.data),
				http.WithService(bootCfg.Name),
				http.WithVersion(bootCfg.Version),
			)...,
		)
	}
}
