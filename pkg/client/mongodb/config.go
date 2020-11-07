package mongodb

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"go.mongodb.org/mongo-driver/event"
)

type (
	Config struct {
		path                 string
		commandMonitor       *event.CommandMonitor
		poolMonitor          *event.PoolMonitor
		URI                  string `config:"uri" desc:"mongodb uri string."`
		EnableCommandMonitor bool   `config:"commandMonitor"`
		EnablePoolMonitor    bool   `config:"poolMonitor"`
	}

	OptionFunc func(*Config)
)

func WithCommandMonitor(monitor *event.CommandMonitor) OptionFunc {
	return func(options *Config) {
		options.commandMonitor = monitor
	}
}

func WithPoolMonitor(monitor *event.PoolMonitor) OptionFunc {
	return func(options *Config) {
		options.poolMonitor = monitor
	}
}

func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)
	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("mongodb client build error: %w", err)
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
