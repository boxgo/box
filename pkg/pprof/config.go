package pprof

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
		Addr string `config:"addr" desc:"pprof http server listen addr"`
	}
)

func StdConfig() *Config {
	cfg := DefaultConfig()

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("pprof build error: %w", err)
	}

	return cfg
}

func DefaultConfig() *Config {
	return &Config{
		Addr: ":9091",
	}
}

func (c *Config) Path() string {
	return "pprof"
}

func (c *Config) Build() *PProf {
	return newPProf(c)
}
