package box_test

import (
	"github.com/boxgo/box/v2"
	"github.com/boxgo/box/v2/logger"
	"github.com/boxgo/box/v2/server/ginserver"
	"github.com/gin-gonic/gin"
)

// Example this is a ping-pong http server.
func Example() {
	ginserver.Default.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, "pong")
	})

	app := box.New(
		box.WithBoxes(
			ginserver.Default,
		),
	)

	if err := app.Run(); err != nil {
		logger.Errorw("server run error: ", "err", err)
	}
}
