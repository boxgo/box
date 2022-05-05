package ginsession

import (
	"fmt"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	// Config 配置
	Config struct {
		path        string
		Redis       string   `config:"redis"`
		CookieName  string   `config:"cookieName" desc:"cookie name"`
		CookieNames []string `config:"cookieNames" desc:"cookie names"`
		KeyPair     string   `config:"keyPair" desc:"cookie value encrypt key pair"`
		KeyPrefix   string   `config:"keyPrefix" desc:"redis save key prefix"`
		MaxLen      int      `config:"maxLen" desc:"max val length"`
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
		logger.Panicf("GinSession load config error: %s", err)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path:       fmt.Sprintf("gin.%s.middlewares.session", key),
		Redis:      "default",
		CookieName: config.ServiceName() + "_sid",
		KeyPrefix:  config.ServiceName() + "_",
		KeyPair:    "",
		MaxLen:     10240,
	}
}

// Build 构建实例
func (c *Config) Build() *GinSession {
	return newGinSession(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
