package main

import (
	"github.com/boxgo/box"
	"github.com/boxgo/box/config"
	"github.com/boxgo/box/config/loader"
	"github.com/boxgo/box/server/rest"
)

func main() {
	app := box.NewBox(
		box.WithConfig(config.NewConfig(loader.NewFileConfig("example/dev.yaml"))),
	)

	app.Mount(&Logger{}, &Schedule{}, rest.NewServer())

	app.Serve()
}
