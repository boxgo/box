// Package schedule is to help you manage schedule tasks.
package schedule

import (
	"context"
	"fmt"
	"net/url"
	"runtime/debug"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/locker"
	"github.com/boxgo/box/pkg/locker/redislocker"
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/metric"
	"github.com/boxgo/box/pkg/server/ginserver"
	"github.com/boxgo/box/pkg/trace"
	"github.com/boxgo/box/pkg/util/strutil"
	"github.com/boxgo/box/pkg/util/urlutil"
	"github.com/robfig/cron"
)

type (
	// Schedule instance is a server, you should mount bo box application to manage lifecycle.
	Schedule struct {
		cfg           *Config
		cron          *cron.Cron
		locker        locker.MutexLocker
		recorder      Recorder
		onceHandler   Handler
		timingHandler Handler
	}

	Handler func(ctx context.Context) error

	Journal struct {
		Config    Config
		StartTime time.Time
		EndTime   time.Time
		Error     error
		Panic     interface{}
	}
	Recorder func(Journal)

	argsKey  struct{}
	queryKey struct{}
)

var (
	scheduleCounter = metric.NewCounterVec(
		"schedule_total",
		"schedule counter",
		[]string{"task", "error", "panic"},
	)
)

func newSchedule(cfg *Config) *Schedule {
	sch := &Schedule{
		cfg:           cfg,
		cron:          cron.New(),
		locker:        cfg.locker,
		recorder:      cfg.recorder,
		onceHandler:   cfg.onceHandler,
		timingHandler: cfg.timingHandler,
	}

	if sch.locker == nil {
		sch.locker = redislocker.Default
	}

	if sch.onceHandler == nil && sch.timingHandler == nil {
		logger.Panic("schedule handler is not set")
	}

	if sch.cfg.Type == Once {
		return sch
	}

	var specs []string
	if len(sch.cfg.Specs) != 0 {
		specs = sch.cfg.Specs
	} else {
		specs = []string{sch.cfg.Spec}
	}

	for _, spec := range specs {
		if err := sch.cron.AddFunc(spec, func() {
			sch.exec(sch.timingHandler)
		}); err != nil {
			logger.Panicf("schedule build error: %s", err)
		}
	}

	if cfg.server != nil {
		if err := sch.serverHttp(); err != nil {
			logger.Panicf("schedule build serverHttp.error: %s", err)
		}
	}

	return sch
}

func (sch *Schedule) Name() string {
	return "schedule"
}

// Serve schedule
func (sch *Schedule) Serve(context.Context) error {
	switch sch.cfg.Type {
	case Once:
		sch.execOnce()
	case Timing:
		sch.execTiming()
	case OnceAndTiming:
		sch.execOnce()
		sch.execTiming()
	}

	return nil
}

// Shutdown stop cron
func (sch *Schedule) Shutdown(context.Context) error {
	if sch.cron != nil {
		sch.cron.Stop()
	}

	return nil
}

// ExecOnce exec once handler immediately
func (sch *Schedule) ExecOnce() {
	sch.execOnce()
}

// ExecTiming exec timing handler immediately
func (sch *Schedule) ExecTiming() {
	if sch.timingHandler == nil {
		return
	}

	sch.exec(sch.timingHandler)
}

func (sch *Schedule) execOnce() {
	if sch.onceHandler == nil {
		return
	}

	sch.exec(sch.onceHandler)
}

func (sch *Schedule) execTiming() {
	if sch.timingHandler == nil {
		return
	}

	sch.cron.Start()
}

func (sch *Schedule) exec(handler Handler) {
	go func() {
		var (
			ctx, cancel = context.WithTimeout(context.TODO(), sch.cfg.Timeout)
			journal     = Journal{
				Config:    *sch.cfg,
				StartTime: time.Now(),
				EndTime:   time.Time{},
				Error:     nil,
				Panic:     nil,
			}
		)

		defer cancel()

		ctx = context.WithValue(ctx, argsKey{}, sch.cfg.Args)
		ctx = context.WithValue(ctx, trace.BizID(), sch.key())
		ctx = context.WithValue(ctx, trace.ReqID(), strutil.RandomAlphanumeric(10))

		if sch.cfg.Delay > 0 {
			time.Sleep(sch.cfg.Delay)
		}

		defer func() {
			if sch.cfg.Compete && sch.cfg.AutoUnlock {
				if err := sch.locker.UnLock(context.Background(), sch.key()); err != nil {
					logger.Trace(ctx).Errorf("Schedule unlock error: [%s]", err)
				} else {
					logger.Trace(ctx).Infof("Schedule unlock success")
				}
			}
		}()

		if sch.cfg.Compete {
			locked, err := sch.locker.Lock(context.Background(), sch.key(), sch.cfg.lockDuration)
			if err != nil {
				logger.Trace(ctx).Errorf("Schedule compete error: [%s]", err)
				return
			} else if !locked {
				logger.Trace(ctx).Warnf("Schedule compete fail")
				return
			} else {
				logger.Trace(ctx).Infof("Schedule compete win")
			}
		}

		// only record winner
		defer func() {
			journal.EndTime = time.Now()
			journal.Panic = recover()

			if journal.Panic != nil {
				logger.Trace(ctx).Errorf("Schedule crash: %+v\n%s", journal.Panic, debug.Stack())
				scheduleCounter.WithLabelValues(sch.key(), "", fmt.Sprintf("%s", journal.Panic)).Inc()
			}

			sch.recorder(journal)
		}()

		logger.Trace(ctx).Infof("Schedule run start")

		if journal.Error = handler(ctx); journal.Error != nil {
			logger.Trace(ctx).Errorf("Schedule run error: [%s]", journal.Error)
			scheduleCounter.WithLabelValues(sch.key(), journal.Error.Error(), "").Inc()
		} else {
			logger.Trace(ctx).Infof("Schedule run success")
			scheduleCounter.WithLabelValues(sch.key(), "", "").Inc()
		}
	}()
}

func (sch Schedule) serverHttp() error {
	endpoint, err := urlutil.UrlJoin("/schedule", sch.cfg.key)
	if err != nil {
		return err
	}

	sch.cfg.server.GET(endpoint, func(ctx *ginserver.Context) {
		var err error
		typ := ctx.DefaultQuery("type", "timing")

		switch c := context.WithValue(ctx, queryKey{}, ctx.Request.URL.Query()); typ {
		case "once":
			err = sch.cfg.onceHandler(c)
		case "timing":
			err = sch.cfg.timingHandler(c)
		default:
			err = sch.cfg.timingHandler(c)
		}

		if err != nil {
			ctx.String(200, "run %s.%sHandler error: %s", sch.cfg.key, typ, err.Error())
		} else {
			ctx.String(500, "run %s.%sHandler success", sch.cfg.key, typ)
		}
	})

	return nil
}

func (sch *Schedule) key() string {
	return fmt.Sprintf("%s.%s", config.ServiceName(), sch.cfg.path)
}

func ArgsVal(ctx context.Context) Args {
	if val, ok := ctx.Value(argsKey{}).(Args); ok {
		return val
	}

	return Args{}
}

func QueryVal(ctx context.Context) url.Values {
	if val, ok := ctx.Value(queryKey{}).(url.Values); ok {
		return val
	}

	return url.Values{}
}
