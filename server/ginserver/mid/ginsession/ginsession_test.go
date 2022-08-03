package ginsession_test

import (
	"github.com/boxgo/box/v2/logger"
	ginserver2 "github.com/boxgo/box/v2/server/ginserver"
	"github.com/boxgo/box/v2/server/ginserver/mid/ginsession"
)

func ExampleGinSession_Redis() {
	ginserver2.Use(ginsession.Redis())
	ginserver2.GET("/ping", func(ctx *ginserver2.Context) {
		ctx.Data(200, "text/plain", []byte("pong"))
	})

	if err := ginserver2.Run(); err != nil {
		logger.Fatal(err)
	}
}

func ExampleGinSession_Cookie() {
	ginserver2.Use(ginsession.Cookie())
	ginserver2.GET("/ping", func(ctx *ginserver2.Context) {
		ctx.Data(200, "text/plain", []byte("pong"))
	})

	if err := ginserver2.Run(); err != nil {
		logger.Fatal(err)
	}
}
