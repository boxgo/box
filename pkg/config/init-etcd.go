// +build !no_config_init,!no_config_etcd

package config

import (
	"github.com/boxgo/box/pkg/config/source/etcd"
)

func init() {
	if !bootOK {
		return
	}

	for idx, cfg := range bootCfg.Source {
		if cfg.name != "etcd" || len(cfg.data) == 0 {
			continue
		}

		defaultSources[idx] = etcd.NewSource(etcd.WithConfig(cfg.data)...)
	}
}
