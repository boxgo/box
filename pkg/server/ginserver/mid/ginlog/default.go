package ginlog

import "github.com/gin-gonic/gin"

var (
	Default = StdConfig("default").Build()
)

func Logger() func(ctx *gin.Context) {
	return Default.Logger()
}
