package mongodb

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
		path                 string
		monitor              Monitor
		URI                  string `config:"uri" desc:"mongodb uri string."`
		EnableCommandMonitor bool   `config:"commandMonitor"`
		EnablePoolMonitor    bool   `config:"poolMonitor"`
	}

	OptionFunc func(*Config)
)

func WithMonitor(monitor Monitor) OptionFunc {
	return func(options *Config) {
		options.monitor = monitor
	}
}

func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)
	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("mongodb client build error: %s", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:                 "mongo." + key,
		URI:                  "mongodb://127.0.0.1:27017",
		EnableCommandMonitor: true,
		EnablePoolMonitor:    false,
	}
}

func (c *Config) Build() *Mongo {
	return newMongo(c)
}

func (c *Config) Path() string {
	return c.path
}
