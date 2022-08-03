package ginprom

import (
	"github.com/gin-gonic/gin"
)

var (
	Default = StdConfig("default").Build()
)

func Handler() gin.HandlerFunc {
	return Default.Handler()
}
