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
)

func DefaultConfig() *Config {
	return &Config{
		Namespace:     "",
		Subsystem:     "",
		PushEnabled:   false,
		PushTargetURL: "",
		PushInterval:  time.Second * 3,
	}
}

func (c *Config) Build() *Metric {
	cfg := DefaultConfig()

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("metric build error: %w", err)
	}

	return newMetric(cfg)
}

func (c *Config) Path() string {
	return "metric"
}
