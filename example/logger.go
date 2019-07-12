package main

import (
	"context"
	"time"

	"github.com/boxgo/logger"
)

type (
	Logger struct {
		Level          string `config:"level" default:"info" help:"Levels: debug,info,warn,error,dpanic,panic,fatal"`
		Encoding       string `config:"encoding" help:"PS: console or json"`
		TraceUID       string `config:"traceUid" help:"Name as trace uid in context"`
		TraceRequestID string `config:"traceRequestId" help:"Name as trace requestId in context"`
	}
)

func (l *Logger) Name() string {
	return "logger"
}

func (l *Logger) Serve(ctx context.Context) error {
	logger.Default.Debug(l.Name(), " Serve")

	for {
		time.Sleep(time.Second)
	}

	return nil
}

func (l *Logger) Shutdown(ctx context.Context) error {
	logger.Default.Debug(l.Name(), " Shutdown")
	return nil
}

func (l *Logger) ConfigWillLoad(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ConfigWillLoad")
}
func (l *Logger) ConfigDidLoad(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ConfigDidLoad")
}

func (l *Logger) ServerWillReady(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ServerWillReady")
}
func (l *Logger) ServerDidReady(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ServerDidReady")
}
func (l *Logger) ServerWillClose(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ServerWillClose")
}
func (l *Logger) ServerDidClose(ctx context.Context) {
	logger.Default.Debug(l.Name(), " ServerDidClose")
}
