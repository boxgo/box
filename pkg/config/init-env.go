//go:build !no_config_init && !no_config_env
// +build !no_config_init,!no_config_env

package config

import (
	"encoding/json"
	"strings"

	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/config/source/env"
)

type (
	envSource struct {
		Prefix string `json:"prefix"`
	}
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "env" || len(cfg.data) == 0 {
			continue
		}

		envSour := envSource{}
		json.Unmarshal(cfg.data, &envSour)

		prefix := strings.ToUpper(strings.ReplaceAll(envSour.Prefix, "-", "_"))

		var opts []source.Option

		if prefix != "" {
			opts = append(opts, env.WithStrippedPrefix(prefix))
		}

		defaultSources[idx] = env.NewSource(opts...)
	}
}
