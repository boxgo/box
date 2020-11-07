package ginserver

import (
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	Config struct {
		path         string
		Mode         string        `config:"mode" desc:"Gin mode: debug,release,test. default is release"`
		Addr         string        `config:"addr"`
		ReadTimeout  time.Duration `config:"readTimeout"`
		WriteTimeout time.Duration `config:"writeTimeout"`
		IdleTimeout  time.Duration `config:"idleTimeout"`
	}
)

func StdConfig(key string) *Config {
	cfg := DefaultConfig(key)

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("gin server build error: %w", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:         "gin." + key,
		Mode:         gin.ReleaseMode,
		Addr:         ":8000",
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute * 5,
	}

}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build() *GinServer {
	return newGinServer(c)
}