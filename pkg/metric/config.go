package metric

import (
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
		Namespace     string        `config:"namespace"`
		Subsystem     string        `config:"subsystem"`
		PushEnabled   bool          `config:"pushEnabled"`
		PushTargetURL string        `config:"pushTargetURL"`
		PushInterval  time.Duration `config:"pushInterval"`
	}

	// OptionFunc is option function.
	OptionFunc func(*Config)
)

// StdConfig load config from config center.
func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)
	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("Metric load config error: %s", err)
	}

	return cfg
}

// DefaultConfig is default config.
func DefaultConfig(key string) *Config {
	return &Config{
		Namespace:     "",
		Subsystem:     "",
		PushEnabled:   false,
		PushTargetURL: "",
		PushInterval:  time.Second * 3,
	}
}

// Build a instance.
func (c *Config) Build() *Metric {
	return newMetric(c)
}

// Path is config path.
func (c *Config) Path() string {
	return "metric"
}
