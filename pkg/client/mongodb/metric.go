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
	Monitor struct {
		stopWatch chan bool
	}
)

var (
	cmdTotal = metric.NewCounterVec(
		"mongo_client_command_total",
		"mongodb client command counter",
		[]string{"command", "error"},
	)
	cmdDuration = metric.NewSummaryVec(
		"mongo_client_command_duration_seconds",
		"mongodb client command duration seconds",
		[]string{"command", "error"},
		map[float64]float64{
			0.25: 0.05,
			0.5:  0.05,
			0.75: 0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	)
	workingSession = metric.NewGaugeVec(
		"mongo_client_session_in_progress",
		"mongo client session in progress gauge",
		[]string{},
	)
	defaultMonitor = newMonitor()
)

func newMonitor() *Monitor {
	m := &Monitor{
		stopWatch: make(chan bool),
	}

	return m
}

func (mon *Monitor) CommandMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started:   mon.started,
		Succeeded: mon.succeeded,
		Failed:    mon.failed,
	}
}

func (mon *Monitor) PoolEventMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{
		Event: mon.event,
	}
}

func (mon *Monitor) started(ctx context.Context, ev *event.CommandStartedEvent) {
	logger.Trace(ctx).Debugf("mongo_command_start cmd: %s, reqId: %d, connId: %s, db: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, ev.DatabaseName)
}

func (mon *Monitor) succeeded(ctx context.Context, ev *event.CommandSucceededEvent) {
	labels := []string{ev.CommandName, ""}
	cmdTotal.WithLabelValues(labels...).Inc()
	cmdDuration.WithLabelValues(labels...).Observe(time.Duration(ev.DurationNanos).Seconds())

	logger.Trace(ctx).Debugf("mongo_command_success cmd: %s, reqId: %d, connId: %s, duration: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, time.Duration(ev.DurationNanos))
}

func (mon *Monitor) failed(ctx context.Context, ev *event.CommandFailedEvent) {
	labels := []string{ev.CommandName, ev.Failure}
	cmdTotal.WithLabelValues(labels...).Inc()
	cmdDuration.WithLabelValues(labels...).Observe(time.Duration(ev.DurationNanos).Seconds())

	logger.Trace(ctx).Debugf("mongo_command_error cmd: %s, reqId: %d, connId: %s, duration: %d, error: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, ev.DurationNanos, ev.Failure)
}

func (mon *Monitor) event(ev *event.PoolEvent) {
	logger.Debugf("mongo_pool_event type: %s, address: %s, connId: %d, reason: %s", ev.Type, ev.Address, ev.ConnectionID, ev.Reason)
}

func (mon *Monitor) watch(db *mongo.Client) {
	go func() {
		for {
			time.Sleep(time.Second)

			select {
			case <-mon.stopWatch:
				logger.Debugf("mongo monitor watch exit")
				close(mon.stopWatch)
				break
			default:
				workingSession.WithLabelValues().Set(float64(db.NumberSessionsInProgress()))
			}
		}
	}()
}

func (mon *Monitor) shutdown() {
	mon.stopWatch <- true
}
