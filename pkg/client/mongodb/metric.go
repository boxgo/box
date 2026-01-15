package mongodb

import (
	"context"
	"time"

	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/metric"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	metricMonitor struct {
		stopWatch chan bool
		client    *mongo.Client
	}
)

var (
	cmdTotal = metric.NewCounterVec(
		"mongo_client_requests_total",
		"The total number of MongoDB commands executed.",
		[]string{"command", "result"},
	)
	cmdDuration = metric.NewHistogramVec(
		"mongo_client_request_duration_seconds",
		"The MongoDB command latencies in seconds.",
		[]string{"command", "result"},
		// 250us, 500us, 1ms, 2.5ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s
		[]float64{0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5},
	)
	workingSession = metric.NewGaugeVec(
		"mongo_client_sessions_inflight",
		"The number of MongoDB sessions currently in progress.",
		[]string{},
	)
)

func newMonitor() Monitor {
	m := &metricMonitor{
		stopWatch: make(chan bool),
	}

	return m
}

func (mon *metricMonitor) Setup(client *mongo.Client) {
	mon.client = client
}

func (mon *metricMonitor) Serve() {
	go func() {
		for {
			time.Sleep(time.Second)

			select {
			case <-mon.stopWatch:
				logger.Debugf("mongo monitor watch exit")
				close(mon.stopWatch)
				break
			default:
				workingSession.WithLabelValues().Set(float64(mon.client.NumberSessionsInProgress()))
			}
		}
	}()
}

func (mon *metricMonitor) Shutdown() {
	mon.stopWatch <- true
}

func (mon *metricMonitor) Started(ctx context.Context, ev *event.CommandStartedEvent) {
	logger.Trace(ctx).Debugf("mongo_command_start cmd: %s, reqId: %d, connId: %s, db: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, ev.DatabaseName)
}

func (mon *metricMonitor) Succeeded(ctx context.Context, ev *event.CommandSucceededEvent) {
	labels := []string{ev.CommandName, "success"}
	cmdTotal.WithLabelValues(labels...).Inc()
	cmdDuration.WithLabelValues(labels...).Observe(time.Duration(ev.DurationNanos).Seconds())

	logger.Trace(ctx).Debugf("mongo_command_success cmd: %s, reqId: %d, connId: %s, duration: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, time.Duration(ev.DurationNanos))
}

func (mon *metricMonitor) Failed(ctx context.Context, ev *event.CommandFailedEvent) {
	labels := []string{ev.CommandName, "error"}
	cmdTotal.WithLabelValues(labels...).Inc()
	cmdDuration.WithLabelValues(labels...).Observe(time.Duration(ev.DurationNanos).Seconds())

	logger.Trace(ctx).Debugf("mongo_command_error cmd: %s, reqId: %d, connId: %s, duration: %d, error: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, ev.DurationNanos, ev.Failure)
}

func (mon *metricMonitor) Event(ev *event.PoolEvent) {
	logger.Debugf("mongo_pool_event type: %s, address: %s, connId: %d, reason: %s", ev.Type, ev.Address, ev.ConnectionID, ev.Reason)
}
