package gopg

import (
	"context"
	"crypto/tls"
	"net"
	"runtime"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"github.com/go-pg/pg/v10"
)

type (
	// Config 配置
	Config struct {
		path                  string
		tLSConfig             *tls.Config
		dialer                func(ctx context.Context, network, addr string) (net.Conn, error)
		onConnect             func(ctx context.Context, cn *pg.Conn) error
		Debug                 bool          `config:"debug" desc:"print all queries (even those without an error)"`
		URI                   string        `config:"uri" desc:"pg connection url. example: postgres://user:pass@localhost:5432/db_name?k=v"`
		ApplicationName       string        `config:"applicationName" desc:"ApplicationName is the application name. Used in logs on Pg side. Only available from pg-9.0."`
		Network               string        `config:"network" desc:"Network type, either tcp or unix. Default is tcp."`
		Addr                  string        `config:"addr" desc:"TCP host:port or Unix socket depending on Network."`
		User                  string        `config:"user"`
		Password              string        `config:"password"`
		Database              string        `config:"database"`
		DialTimeout           time.Duration `config:"dialTimeout" desc:"Dial timeout for establishing new connections. Default is 5 seconds."`
		ReadTimeout           time.Duration `config:"readTimeout" desc:"Timeout for socket reads. If reached, commands will fail and user is authenticated."`
		WriteTimeout          time.Duration `config:"writeTimeout" desc:"Timeout for socket writes. If reached, commands will fail with a timeout instead of blocking."`
		MaxRetries            int           `config:"maxRetries" desc:"Maximum number of retries before giving up. Default is to not retry failed queries."`
		RetryStatementTimeout bool          `config:"retryStatementTimeout" desc:"Whether to retry queries cancelled because of statement_timeout."`
		MinRetryBackoff       time.Duration `config:"minRetryBackoff" desc:"Minimum backoff between each retry. Default is 250 milliseconds; -1 disables backoff."`
		MaxRetryBackoff       time.Duration `config:"maxRetryBackoff" desc:"Maximum backoff between each retry. Default is 4 seconds; -1 disables backoff."`
		PoolSize              int           `config:"poolSize" desc:"Maximum number of socket connections. Default is 10 connections per every CPU as reported by runtime.NumCPU."`
		MinIdleConns          int           `config:"minIdleConns" desc:"Minimum number of idle connections which is useful when establishing new connection is slow. Default is 1."`
		MaxConnAge            time.Duration `config:"maxConnAge" desc:"Connection age at which client retires (closes) the connection. It is useful with proxies like PgBouncer and HAProxy. Default is to not close aged connections."`
		PoolTimeout           time.Duration `config:"poolTimeout" desc:"Time for which client waits for free connection if all connections are busy before returning an error. Default is 30 seconds if ReadTimeOut is not defined, otherwise, ReadTimeout + 1 second."`
		IdleTimeout           time.Duration `config:"idleTimeout" desc:"Amount of time after which client closes idle connections. Should be less than server's timeout. Default is 5 minutes. -1 disables idle timeout check."`
		IdleCheckFrequency    time.Duration `config:"idleCheckFrequency" desc:"Frequency of idle checks made by idle connections reaper. Default is 1 minute. -1 disables idle connections reaper, but idle connections are still discarded by the client if IdleTimeout is set."`
	}

	// OptionFunc 选项信息
	OptionFunc func(*Config)
)

func WithTLSConfig(cfg *tls.Config) OptionFunc {
	return func(c *Config) {
		c.tLSConfig = cfg
	}
}

func WithDialer(fn func(ctx context.Context, network, addr string) (net.Conn, error)) OptionFunc {
	return func(c *Config) {
		c.dialer = fn
	}
}

func WithOnConnect(fn func(ctx context.Context, cn *pg.Conn) error) OptionFunc {
	return func(c *Config) {
		c.onConnect = fn
	}
}

// StdConfig 标准配置
func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)
	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("PostgreSQL load config error: %s", err)
	} else {
		logger.Debugw("PostgreSQL load config", "config", cfg)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	return &Config{
		path:                  "pg." + key,
		Debug:                 false,
		URI:                   "",
		tLSConfig:             nil,
		dialer:                nil,
		onConnect:             nil,
		ApplicationName:       config.ServiceName(),
		Network:               "tcp",
		Addr:                  "",
		User:                  "",
		Password:              "",
		Database:              "",
		DialTimeout:           time.Second * 5,
		ReadTimeout:           time.Second * 5,
		WriteTimeout:          time.Second * 5,
		MaxRetries:            0,
		RetryStatementTimeout: false,
		MinRetryBackoff:       250 * time.Millisecond,
		MaxRetryBackoff:       4 * time.Second,
		PoolSize:              10 * runtime.NumCPU(),
		MinIdleConns:          1,
		MaxConnAge:            time.Hour * 1,
		PoolTimeout:           time.Second * 10,
		IdleTimeout:           time.Minute * 5,
		IdleCheckFrequency:    time.Minute * 1,
	}
}

// Build 构建实例
func (c *Config) Build() *PostgreSQL {
	return newPostgreSQL(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
