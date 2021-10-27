package redis

import (
	"context"
	"log"
	"time"

	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/util/strutil"
	"github.com/go-redis/redis/v8"
)

type (
	redisSource struct {
		err    error
		prefix string
		opts   source.Options
		client redis.UniversalClient
	}
)

func NewSource(opts ...source.Option) source.Source {
	var (
		options = source.NewOptions(opts...)
		prefix  = "box"
		client  redis.UniversalClient
	)

	if val, ok := options.Context.Value(redisConfigKey{}).(redisConfig); !ok {
		log.Panic("config source redis is not set.")
	} else {
		client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        strutil.Compact(val.Address),
			DB:           val.DB,
			Password:     val.Password,
			PoolSize:     val.PoolSize,
			MasterName:   val.MasterName,
			MinIdleConns: val.MinIdleConnCnt,
		})
	}

	if val, ok := options.Context.Value(prefixKey{}).(string); ok && val != "" {
		prefix = val
	}

	return &redisSource{
		err:    nil,
		prefix: prefix,
		client: client,
		opts:   options,
	}
}

func (rs *redisSource) Read() (*source.ChangeSet, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	rsp, err := rs.client.Get(context.Background(), rs.prefix+".config").Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Source:    rs.String(),
		Data:      rsp,
		Format:    rs.opts.Encoder.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil

}
func (rs *redisSource) Watch() (source.Watcher, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	return newWatcher(rs)
}

func (rs *redisSource) String() string {
	return "redis"
}
