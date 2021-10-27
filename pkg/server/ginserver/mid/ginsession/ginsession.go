// Package ginsession is gin server session middleware.
package ginsession

import (
	"github.com/boxgo/box/pkg/client/redis"
	"github.com/boxgo/redisstore/v2"
	"github.com/boxgo/redisstore/v2/serializer"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type (
	GinSession struct {
		cfg    *Config
		client *redis.Redis
	}
)

func newGinSession(c *Config) *GinSession {
	client := redis.StdConfig(c.Redis).Build()

	return &GinSession{
		cfg:    c,
		client: client,
	}

}

func (s *GinSession) Cookie() gin.HandlerFunc {
	return sessions.Sessions(s.cfg.CookieName, cookie.NewStore([]byte(s.cfg.KeyPair)))
}

func (s *GinSession) Redis() gin.HandlerFunc {
	st, _ := redisstore.NewStoreWithUniversalClient(
		s.client.Client(),
		redisstore.WithMaxLength(s.cfg.MaxLen),
		redisstore.WithKeyPrefix(s.cfg.KeyPrefix),
		redisstore.WithKeyPairs([]byte(s.cfg.KeyPair)),
		redisstore.WithSerializer(serializer.JSONSerializer{}),
	)

	return sessions.Sessions(s.cfg.CookieName, &redisStore{
		RedisStore: st,
	})
}
