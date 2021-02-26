package gopg

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	// Config 配置
	Config struct {
		path  string
		Debug bool   `config:"debug" desc:"print all queries (even those without an error)"`
		URI   string `config:"uri" desc:"pg connection url. example: postgres://user:pass@localhost:5432/db_name?k=v"`
	}

	// OptionFunc 选项信息
	OptionFunc func(*Config)
)

// StdConfig 标准配置
func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)
	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("PostgreSQL build error: %s", err)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path:  "pg." + key,
		Debug: false,
		URI:   "postgres://user:pass@localhost:5432/db_name",
	}
}

// Build 构建实例
func (c *Config) Build() *PostgreSQL {
	return newPostgreSQL(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
