package main

import (
	"context"

	"github.com/boxgo/logger"
)

type (
	Schedule struct {
		Server string `config:"server" desc:"server ip"`
	}
)

func (l *Schedule) Name() string {
	return "schedule"
}

func (l *Schedule) Serve(ctx context.Context) error {
	logger.Default.Debug(l.Name(), " Serve")
	return nil
}

func (l *Schedule) Shutdown(ctx context.Context) error {
	logger.Default.Debug(l.Name(), " Shutdown")
	return nil
}

func (l *Schedule) ConfigWillLoad(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ConfigWillLoad")
}

func (l *Schedule) ConfigDidLoad(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ConfigDidLoad")
}

func (l *Schedule) ServerWillReady(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ServerWillReady")
}

func (l *Schedule) ServerDidReady(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ServerDidReady")
}

func (l *Schedule) ServerWillClose(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ServerWillClose")
}

func (l *Schedule) ServerDidClose(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ServerDidClose")
}
