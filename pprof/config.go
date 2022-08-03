package pprof

import (
	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
)

type (
	Config struct {
	}
)

func StdConfig() *Config {
	cfg := DefaultConfig()

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("pprof build error: %s", err)
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
