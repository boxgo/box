package ginprom

import (
	"fmt"
	"strings"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	// Config 配置
	Config struct {
		path                string
		requestURLMappingFn func(*gin.Context) string
	}

	// OptionFunc 选项信息
	OptionFunc func(*Config)
)

// WithURLMapping set up url mapping func
// default is: UrlMapping
func WithURLMapping(fn func(*gin.Context) string) OptionFunc {
	return func(options *Config) {
		options.requestURLMappingFn = fn
	}
}

// StdConfig 标准配置
func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)
	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("GinProm load config error: %s", err)
	} else {
		logger.Debugw("GinProm load config", "config", cfg)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path:                fmt.Sprintf("gin.%s.middlewares.prom", key),
		requestURLMappingFn: UrlMapping,
	}
}

// Build 构建实例
func (c *Config) Build() *GinProm {
	return newGinProm(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}

func UrlMapping(c *gin.Context) string {
	url := c.Request.URL.Path

	for _, p := range c.Params {
		url = strings.Replace(url, "/"+p.Value, "/:"+p.Key, 1)
	}

	if len(url) > 200 {
		return url[:197] + "..."
	}

	return url
}
