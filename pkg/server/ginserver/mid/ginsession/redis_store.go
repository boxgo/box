package ginsession

import (
	"github.com/boxgo/redisstore/v2"
	"github.com/gin-contrib/sessions"
)

type (
	redisStore struct {
		*redisstore.RedisStore
	}
)

func (st *redisStore) Options(opts sessions.Options) {
	st.RedisStore.SetOptions(opts.ToGorillaOptions())
}
