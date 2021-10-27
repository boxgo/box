package ginsession

import (
	"github.com/gin-gonic/gin"
)

var (
	Default = StdConfig("default").Build()
)

func Cookie() gin.HandlerFunc {
	return Default.Cookie()
}

func Redis() gin.HandlerFunc {
	return Default.Redis()
}
