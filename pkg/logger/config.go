package logger

import (
	"github.com/boxgo/box/pkg/config"
)

type (
	Config struct {
		config.SubConfigurator
		Level    *config.Field
		Encoding *config.Field
	}
)

func newConfig(name string, cfg config.SubConfigurator) *Config {
	c := &Config{
		SubConfigurator: cfg,
		Level:           config.NewField(name, "level", "levels: debug,info,warn,error,dpanic,panic,fatal", "info"),
		Encoding:        config.NewField(name, "encoding", "console or json", "console"),
	}

	c.Mount(c.Fields()...)

	return c
}

func (c *Config) GetLevel() string {
	return c.SubConfigurator.GetString(c.Level)
}

func (c *Config) GetEncoding() string {
	return c.SubConfigurator.GetString(c.Encoding)
}

func (c *Config) Fields() []*config.Field {
	return []*config.Field{
		c.Level,
		c.Encoding,
	}
}
