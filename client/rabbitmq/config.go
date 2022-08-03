package rabbitmq

import (
	"time"

	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
)

type (
	// Config 配置
	Config struct {
		path       string
		URI        string        `config:"uri" desc:"Connection uri"`
		Vhost      string        `config:"vhost" desc:"Vhost specifies the namespace of permissions, exchanges, queues and bindings on the server. Dial sets this to the path parsed from the URL."`
		ChannelMax int           `config:"channelMax" desc:"0 max channels means 2^16 - 1"`
		FrameSize  int           `config:"frameSize" desc:"0 max bytes means unlimited"`
		Heartbeat  time.Duration `config:"heartbeat" desc:"less than 1s uses the server's interval"`
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
		logger.Panicf("RabbitMQ load config error: %s", err)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path: "rabbitmq." + key,
	}
}

// Build 构建实例
func (c *Config) Build() *RabbitMQ {
	return newRabbitMQ(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
