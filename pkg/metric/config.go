package metric

import (
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
		Namespace     string            `config:"namespace"`
		Subsystem     string            `config:"subsystem"`
		PushEnabled   bool              `config:"pushEnabled"`
		PushTargetURL string            `config:"pushTargetURL"`
		PushInterval  time.Duration     `config:"pushInterval"`
		PushBasicAuth PushBasicAuth     `config:"pushBasicAuth"`
		BasicAuth     map[string]string `config:"basicAuth" desc:"Basic Auth accounts for metrics. Keys are usernames and values are passwords."`
	}

	// PushBasicAuth holds the credentials sent to the Pushgateway.
	PushBasicAuth struct {
		Username string `config:"username"`
		Password string `config:"password"`
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
	} else {
		logCfg := *cfg
		logCfg.BasicAuth = nil
		logCfg.PushBasicAuth = PushBasicAuth{}
		logger.Debugw("Metric load config", "config", &logCfg)
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
	return &Metric{cfg: c}
}

// Path is config path.
func (c *Config) Path() string {
	return "metric"
}
