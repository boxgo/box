package main

import (
	"github.com/boxgo/box"
	"github.com/boxgo/box/config"
	"github.com/boxgo/config/source/env"
	"github.com/boxgo/config/source/file"
	"github.com/boxgo/config/source/pflag"
)

func main() {
	cfg := config.NewConfig(
		file.NewSource(file.WithPath("./example/dev.yaml")),
		pflag.NewSource(),
		env.NewSource(),
	)

	cfg.EnablePflag(true)

	app := box.NewBox(
		box.WithConfig(cfg),
	)

	app.Mount(&Logger{}, &Schedule{})

	app.Serve()
}
