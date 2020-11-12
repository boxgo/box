package pprof

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
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
	return &Config{}
}

func (c *Config) Path() string {
	return "pprof"
}

func (c *Config) Build() *PProf {
	return newPProf(c)
}
