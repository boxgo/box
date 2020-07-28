package mongodb

import (
	"context"

	"github.com/boxgo/box/pkg/logger"
	"go.mongodb.org/mongo-driver/event"
)

func newCommandMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started:   Started,
		Succeeded: Succeeded,
		Failed:    Failed,
	}
}

func newPoolMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{
		Event: Event,
	}
}

func Started(ctx context.Context, ev *event.CommandStartedEvent) {
	logger.Trace(ctx).Info("Started", ev)
}

func Succeeded(ctx context.Context, ev *event.CommandSucceededEvent) {
	logger.Trace(ctx).Info("Succeeded", ev)
}

func Failed(ctx context.Context, ev *event.CommandFailedEvent) {
	logger.Trace(ctx).Info("Failed", ev)
}

func Event(ev *event.PoolEvent) {
	logger.Info("Event", ev)
}
