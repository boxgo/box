package main

import (
	"context"
	"time"

	"github.com/boxgo/kit/logger"

	"github.com/boxgo/box/minibox"
)

type (
	Logger struct {
		Level          string `json:"level" desc:"Levels: debug,info,warn,error,dpanic,panic,fatal"`
		Encoding       string `json:"encoding" desc:"PS: console or json"`
		TraceUID       string `json:"traceUid" desc:"Name as trace uid in context"`
		TraceRequestID string `json:"traceRequestId" desc:"Name as trace requestId in context"`
	}
)

func (l *Logger) Name() string {
	return "logger"
}

func (l *Logger) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{}
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

// func (l *Logger) ConfigWillLoad(ctx context.Context) {
// 	logger.Default.Debug(l.Name(), " ConfigWillLoad")
// }
// func (l *Logger) ConfigDidLoad(ctx context.Context) {
// 	logger.Default.Debug(l.Name(), " ConfigDidLoad")
// }

func (l *Logger) ServerWillReady(ctx context.Context) {
	logger.Default.Debug(l.Name(), "  ServerWillReady")
}
func (l *Logger) ServerDidReady(ctx context.Context) {
	logger.Default.Debug(l.Name(), "  ServerDidReady")
}
func (l *Logger) ServerWillClose(ctx context.Context) {
	logger.Default.Debug(l.Name(), "  ServerWillClose")
}
func (l *Logger) ServerDidClose(ctx context.Context) {
	logger.Default.Debug(l.Name(), "  ServerDidClose")
}
