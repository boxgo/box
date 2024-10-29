//go:build !no_config_init && !no_config_http
// +build !no_config_init,!no_config_http

package config

import (
	"encoding/json"

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

		httpCfg := http.Config{}
		if err := json.Unmarshal(cfg.data, &httpCfg); err != nil {
			panic(err)
		}

		namespace := bootCfg.Namespace
		service := bootCfg.Service
		version := bootCfg.Version
		if httpCfg.Namespace != "" {
			namespace = httpCfg.Namespace
		}
		if httpCfg.Service != "" {
			service = httpCfg.Service
		}
		if httpCfg.Version != "" {
			version = httpCfg.Version
		}

		defaultSources[idx] = http.NewSource(
			append(
				http.WithConfig(cfg.data),
				http.WithNamespace(namespace),
				http.WithService(service),
				http.WithVersion(version),
			)...,
		)
	}
}
