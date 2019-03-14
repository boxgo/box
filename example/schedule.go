package main

import (
	"context"

	"github.com/boxgo/kit/logger"

	"github.com/boxgo/box/minibox"
)

type (
	App struct {
		AppName string `json:"name" desc:"Application name"`
	}
)

func (c *App) Name() string {
	return "ext"
}

type (
	Info struct {
		App    App
		Common Common
		Server string `json:"server" desc:"server ip"`
	}
)

func (l *Info) Name() string {
	return "schedule"
}

func (l *Info) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{&l.App, &l.Common}
}

func (l *Info) Serve(ctx context.Context) error {
	logger.Default.Debug(l.Name(), " Serve")
	return nil
}

func (l *Info) Shutdown(ctx context.Context) error {
	logger.Default.Debug(l.Name(), " Shutdown")
	return nil
}

// func (l *Info) ConfigWillLoad(ctx context.Context) {
// 	logger.Default.Debug(l.Name(), " ConfigWillLoad")
// }

// func (l *Info) ConfigDidLoad(ctx context.Context) {
// 	logger.Default.Debug(l.Name(), " ConfigDidLoad")
// }

func (l *Info) ServerWillReady(ctx context.Context) {
	logger.Default.Debug(l.Name(), "  ServerWillReady")
}

func (l *Info) ServerDidReady(ctx context.Context) {
	logger.Default.Debug(l.Name(), "  ServerDidReady")
}

func (l *Info) ServerWillClose(ctx context.Context) {
	logger.Default.Debug(l.Name(), "  ServerWillClose")
}

func (l *Info) ServerDidClose(ctx context.Context) {
	logger.Default.Debug(l.Name(), "  ServerDidClose")
}
