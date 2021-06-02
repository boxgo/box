package gorm

import (
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"gorm.io/gorm"
)

type (
	// Config 配置
	Config struct {
		path         string
		dial         gorm.Dialector
		Driver       string        `config:"driver" desc:"SQLDrivers: https://github.com/golang/go/wiki/SQLDrivers"`
		DSN          string        `config:"dsn" desc:"Data Source Name"`
		MaxIdleTime  time.Duration `config:"maxIdleTime" desc:"SetConnMaxIdleTime sets the maximum amount of time a connection may be idle. If d <= 0, connections are not closed due to a connection's idle time."`
		MaxLifeTime  time.Duration `config:"maxLifeTime" desc:"SetConnMaxLifetime sets the maximum amount of time a connection may be reused. If d <= 0, connections are not closed due to a connection's age."`
		MaxIdleConns int           `config:"maxIdleConns" desc:"MaxIdleConns sets the maximum number of connections in the idle connection pool. If MaxOpenConns is greater than 0 but less than the new MaxIdleConns, then the new MaxIdleConns will be reduced to match the MaxOpenConns limit. If n <= 0, no idle connections are retained."`
		MaxOpenConns int           `config:"maxOpenConns" desc:"MaxOpenConns sets the maximum number of open connections to the database. If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than MaxIdleConns, then MaxIdleConns will be reduced to match the new MaxOpenConns limit. If n <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited)."`
	}

	// OptionFunc 选项信息
	OptionFunc func(*Config)
)

func WithDialector(dial gorm.Dialector) OptionFunc {
	return func(c *Config) {
		c.dial = dial
	}
}

// StdConfig 标准配置
func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)
	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("Gorm load config error: %s", err)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path:         "gorm." + key,
		MaxIdleTime:  time.Minute * 5,
		MaxLifeTime:  time.Hour,
		MaxIdleConns: 2,
		MaxOpenConns: 10,
	}
}

// Build 构建实例
func (c *Config) Build() *Gorm {
	return newGorm(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
