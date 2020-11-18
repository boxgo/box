package trace

import (
	"fmt"

	"github.com/boxgo/box/pkg/config"
)

type (
	// Config of system
	Config struct {
		TraceUID    string `config:"traceUid"`
		TraceReqID  string `config:"traceReqId"`
		TraceSpanID string `config:"traceSpanId"`
		TraceBizID  string `config:"traceBizId"`
	}
)

// StdConfig get from config
func StdConfig() *Config {
	cfg := DefaultConfig()

	if err := config.Scan(cfg); err != nil {
		panic(fmt.Errorf("trace build error: %s", err))
	}

	return cfg
}

// DefaultConfig of trace
func DefaultConfig() *Config {
	return &Config{
		TraceUID:    "box.trace.uid",
		TraceReqID:  "box.trace.reqId",
		TraceSpanID: "box.trace.spanId",
		TraceBizID:  "box.trace.bizId",
	}
}

func (c *Config) Path() string {
	return "trace"
}

// Build Trace
func (c *Config) Build() *Trace {
	return newTrace(c)
}
