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
	if len(s.cfg.CookieNames) != 0 {
		return sessions.SessionsMany(s.cfg.CookieNames, cookie.NewStore([]byte(s.cfg.KeyPair)))
	}

	return sessions.Sessions(s.cfg.CookieName, cookie.NewStore([]byte(s.cfg.KeyPair)))
}

func (s *GinSession) Redis(serializers ...serializer.SessionSerializer) gin.HandlerFunc {
	var ser serializer.SessionSerializer

	if len(serializers) == 0 {
		ser = serializer.JSONSerializer{}
	} else {
		ser = serializers[0]
	}

	st, _ := redisstore.NewStoreWithUniversalClient(
		s.client.Client(),
		redisstore.WithMaxLength(s.cfg.MaxLen),
		redisstore.WithKeyPrefix(s.cfg.KeyPrefix),
		redisstore.WithKeyPairs([]byte(s.cfg.KeyPair)),
		redisstore.WithSerializer(ser),
	)

	if len(s.cfg.CookieNames) != 0 {
		return sessions.SessionsMany(s.cfg.CookieNames, &redisStore{
			RedisStore: st,
		})
	}

	return sessions.Sessions(s.cfg.CookieName, &redisStore{
		RedisStore: st,
	})
}

func (s *GinSession) CookieName() string {
	return s.cfg.CookieName
}

func (s *GinSession) CookieNames() []string {
	return s.cfg.CookieNames
}
