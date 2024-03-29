package grpcclient

import (
	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
	"google.golang.org/grpc"
)

type (
	Config struct {
		path                    string
		dialOptions             []grpc.DialOption
		unaryClientInterceptor  []grpc.UnaryClientInterceptor
		streamClientInterceptor []grpc.StreamClientInterceptor
		Dial                    bool   `config:"dial" desc:"dial when build new client."`
		Target                  string `config:"target" desc:"target server."`
	}

	OptionFunc func(*Config)
)

func WithDialOption(opt ...grpc.DialOption) OptionFunc {
	return func(c *Config) {
		c.dialOptions = append(c.dialOptions, opt...)
	}
}

func WithUnaryClientInterceptor(interceptor ...grpc.UnaryClientInterceptor) OptionFunc {
	return func(c *Config) {
		c.unaryClientInterceptor = append(c.unaryClientInterceptor, interceptor...)
	}
}

func WithStreamClientInterceptor(interceptor ...grpc.StreamClientInterceptor) OptionFunc {
	return func(c *Config) {
		c.streamClientInterceptor = append(c.streamClientInterceptor, interceptor...)
	}
}

func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)

	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("gRPC client build error: %s", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:   "grpc_client." + key,
		Dial:   false,
		Target: ":9001",
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build() *Client {
	return newGRpcClient(c)
}
