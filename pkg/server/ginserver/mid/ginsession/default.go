package ginsession

import (
	"github.com/boxgo/redisstore/v2/serializer"
	"github.com/gin-gonic/gin"
)

var (
	Default = StdConfig("default").Build()
)

func Cookie() gin.HandlerFunc {
	return Default.Cookie()
}

func Redis(serializers ...serializer.SessionSerializer) gin.HandlerFunc {
	return Default.Redis(serializers...)
}

func CookieName() string {
	return Default.CookieName()
}

func CookieNames() []string {
	return Default.CookieNames()
}
