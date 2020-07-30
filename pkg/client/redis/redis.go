package redis

import (
	"context"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/dummybox"
	"github.com/go-redis/redis/v8"
)

type (
	Redis struct {
		dummybox.DummyBox
		name           string
		client         redis.UniversalClient
		cfg            config.SubConfigurator
		masterName     *config.Field
		address        *config.Field
		password       *config.Field
		db             *config.Field
		poolSize       *config.Field
		minIdleConnCnt *config.Field
	}

	Options struct {
		name string
		cfg  config.SubConfigurator
	}

	OptionFunc func(*Options)
)

var (
	Default = New()
)

func New(optionFunc ...OptionFunc) *Redis {
	opts := &Options{}
	for _, fn := range optionFunc {
		fn(opts)
	}

	if opts.name == "" {
		opts.name = "redis.default"
	} else {
		opts.name = "redis." + opts.name
	}
	if opts.cfg == nil {
		opts.cfg = config.Default
	}

	r := &Redis{
		name:           opts.name,
		cfg:            opts.cfg,
		masterName:     config.NewField(opts.name, "masterName", "The sentinel master name. Only failover clients.", ""),
		address:        config.NewField(opts.name, "address", "Either a single address or a seed list of host:port addresses of cluster/sentinel nodes.", []string{}),
		password:       config.NewField(opts.name, "password", "Redis password", ""),
		db:             config.NewField(opts.name, "db", "Database to be selected after connecting to the server. Only single-node and failover clients.", 0),
		poolSize:       config.NewField(opts.name, "poolSize", "Connection pool size", 100),
		minIdleConnCnt: config.NewField(opts.name, "minIdleConnCnt", "Min idle connections.", 50),
	}

	opts.cfg.Mount(r.masterName, r.address, r.password, r.db, r.poolSize, r.minIdleConnCnt)

	return r
}

func (r *Redis) Name() string {
	return r.name
}

func (r *Redis) Serve(ctx context.Context) error {
	r.client = redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName:   r.cfg.GetString(r.masterName),
		Addrs:        r.cfg.GetStringSlice(r.address),
		Password:     r.cfg.GetString(r.password),
		DB:           r.cfg.GetInt(r.db),
		PoolSize:     r.cfg.GetInt(r.poolSize),
		MinIdleConns: r.cfg.GetInt(r.minIdleConnCnt),
	})
	r.client.AddHook(&Metric{
		addr:   r.cfg.GetStringSlice(r.address),
		master: r.cfg.GetString(r.masterName),
		db:     r.cfg.GetInt(r.db),
	})

	return r.client.Ping(ctx).Err()
}

func (r *Redis) Shutdown(ctx context.Context) error {
	if r.client != nil {
		return r.client.Close()
	}

	return nil
}
