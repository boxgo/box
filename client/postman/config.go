package postman

import (
	"time"

	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
)

type (
	// Config 配置
	Config struct {
		path     string
		SSL      bool          `config:"ssl"`
		Timeout  time.Duration `config:"timeout"`
		Address  string        `config:"address"`
		PoolSize int           `config:"poolSize"`
		Identity string        `config:"identity"`
		Username string        `config:"username"`
		Password string        `config:"password"`
		Host     string        `config:"host"`
		From     string        `config:"from"`
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
		logger.Panicf("PostMan load config error: %s", err)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path:     "postman." + key,
		Timeout:  time.Second * 30,
		PoolSize: 10,
	}
}

// Build 构建实例
func (c *Config) Build() *PostMan {
	return newPostMan(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
