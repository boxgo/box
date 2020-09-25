package logger

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/config/field"
)

type (
	Config struct {
		config.SubConfigurator
		Level    *field.Field
		Encoding *field.Field
	}
)

func newConfig(name string, cfg config.SubConfigurator) *Config {
	c := &Config{
		SubConfigurator: cfg,
		Level:           field.New(false, "logger", "level", "levels: debug,info,warn,error,dpanic,panic,fatal", "info"),
		Encoding:        field.New(false, "logger", "encoding", "console or json", "console"),
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

func (c *Config) Fields() []*field.Field {
	return []*field.Field{
		c.Level,
		c.Encoding,
	}
}
