package ginprom_test

import (
	"github.com/boxgo/box/v2/logger"
	"github.com/boxgo/box/v2/server/ginserver"
	"github.com/boxgo/box/v2/server/ginserver/mid/ginprom"
	"github.com/gin-gonic/gin"
)

func Example() {
	ginserver.Use(ginprom.Handler())
	ginserver.GET("/ping", func(ctx *gin.Context) {
		ctx.Data(200, "text/plain", []byte("pong"))
	})

	if err := ginserver.Run(); err != nil {
		logger.Fatal(err)
	}
}
