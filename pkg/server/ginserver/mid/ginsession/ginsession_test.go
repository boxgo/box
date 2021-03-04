package ginsession_test

import (
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/server/ginserver"
	"github.com/boxgo/box/pkg/server/ginserver/mid/ginsession"
)

func ExampleGinSession_Redis() {
	ginserver.Use(ginsession.Redis())
	ginserver.GET("/ping", func(ctx *ginserver.Context) {
		ctx.Data(200, "text/plain", []byte("pong"))
	})

	if err := ginserver.Run(); err != nil {
		logger.Fatal(err)
	}
}

func ExampleGinSession_Cookie() {
	ginserver.Use(ginsession.Cookie())
	ginserver.GET("/ping", func(ctx *ginserver.Context) {
		ctx.Data(200, "text/plain", []byte("pong"))
	})

	if err := ginserver.Run(); err != nil {
		logger.Fatal(err)
	}
}
