package ginserver

import (
	"time"

	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
	"github.com/gin-gonic/gin"
)

type (
	Config struct {
		path         string
		Mode         string            `config:"mode" desc:"Gin mode: debug,release,test. default is release"`
		Addr         string            `config:"addr" desc:"server listen addr, format is ip:port"`
		BasicAuth    map[string]string `config:"basicAuth" desc:"basicAuth. key is username, value is password"`
		ReadTimeout  time.Duration     `config:"readTimeout"`
		WriteTimeout time.Duration     `config:"writeTimeout"`
		IdleTimeout  time.Duration     `config:"idleTimeout"`
	}
)

func StdConfig(key string) *Config {
	cfg := DefaultConfig(key)

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("gin server build error: %s", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	var addr string
	if key == "insight" {
		addr = ":9999"
	} else {
		addr = ":9000"
	}

	return &Config{
		path:         "gin." + key,
		Mode:         gin.ReleaseMode,
		Addr:         addr,
		BasicAuth:    nil,
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
