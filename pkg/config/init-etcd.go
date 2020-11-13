// +build !no_config_init,!no_config_etcd

package config

import (
	"github.com/boxgo/box/pkg/config/source/etcd"
)

func init() {
	if !bootOK {
		return
	}

	for _, cfg := range sourceConfigs {
		if cfg.name != "etcd" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[cfg.idx] = etcd.NewSource(etcd.WithConfig(cfg.data)...)
	}
}
