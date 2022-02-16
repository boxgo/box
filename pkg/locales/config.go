package locales

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	// Config 配置
	Config struct {
		path        string
		LanguageDir string `config:"languageDir" desc:"language files directory"`
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
		logger.Panicf("Locales load config error: %s", err)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path:        "locales." + key,
		LanguageDir: "./languages",
	}
}

// Build 构建实例
func (c *Config) Build() *Locales {
	return newLocales(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
