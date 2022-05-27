package schedule

import (
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/locker"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
		path          string
		onceHandler   Handler
		timingHandler Handler
		lockDuration  time.Duration
		locker        locker.MutexLocker
		recorder      Recorder
		Type          Type                   `config:"type" desc:"Stop: 0, Once: 1, Timing: 2, OnceAndTiming: 3"`
		Spec          string                 `config:"spec" desc:"Cron spec info"`
		Specs         []string               `config:"specs" desc:"Multi cron spec info, higher priority than spec"`
		Compete       bool                   `config:"compete" desc:"Only winner can exec schedule"`
		AutoUnlock    bool                   `config:"autoUnlock" desc:"Auto unlock after task finish"`
		LockSeconds   int                    `config:"lockSeconds" desc:"Lock ttl"`
		Delay         time.Duration          `config:"delay" desc:"Delay duration"`
		Args          map[string]interface{} `config:"args" desc:"Schedule arguments"`
	}

	Type int

	OptionFunc func(*Config)
)

const (
	Stop          = Type(0)
	Once          = Type(1)
	Timing        = Type(2)
	OnceAndTiming = Type(3)
)

func WithRecorder(recorder Recorder) OptionFunc {
	return func(c *Config) {
		c.recorder = recorder
	}
}

func WithLocker(locker locker.MutexLocker) OptionFunc {
	return func(c *Config) {
		c.locker = locker
	}
}

func WithHandler(onceHandler, timingHandler Handler) OptionFunc {
	return func(c *Config) {
		c.onceHandler = onceHandler
		c.timingHandler = timingHandler
	}
}

func StdConfig(key string) *Config {
	cfg := DefaultConfig(key)

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("schedule build error: %s", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:        "schedule." + key,
		Type:        Stop,
		Spec:        "",
		Specs:       []string{},
		Compete:     true,
		AutoUnlock:  true,
		LockSeconds: 0,
		Delay:       0,
		Args:        map[string]interface{}{},
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build(optionFunc ...OptionFunc) *Schedule {
	for _, fn := range optionFunc {
		fn(c)
	}

	c.lockDuration = time.Duration(1000000000 * c.LockSeconds)

	if c.recorder == nil {
		c.recorder = func(Journal) {}
	}

	return newSchedule(c)
}
