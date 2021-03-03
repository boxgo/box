package ginlog

import (
	"fmt"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	// Config 配置
	Config struct {
		path          string
		UserAgent     bool `config:"useragent" desc:"log user-agent"`
		RequestIP     bool `config:"requestIp" desc:"log request ip"`
		RequestHeader bool `config:"requestHeader" desc:"log request header"`
		RequestQuery  bool `config:"requestQuery" desc:"log request query"`
		RequestBody   bool `config:"requestBody" desc:"log request body"`
		ResponseBody  bool `config:"responseBody" desc:"log response body"`
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
		logger.Panicf("GinLog load config error: %s", err)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path:          fmt.Sprintf("gin.%s.middlewares.logger", key),
		UserAgent:     true,
		RequestIP:     true,
		RequestHeader: false,
		RequestQuery:  true,
		RequestBody:   true,
		ResponseBody:  true,
	}
}

// Build 构建实例
func (c *Config) Build() *GinLog {
	return newGinLog(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
