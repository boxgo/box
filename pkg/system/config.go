package system

import (
	"github.com/boxgo/box/pkg/config"
)

type (
	// Config of system
	Config struct {
		Name        string `config:"name"`
		Version     string `config:"version"`
		TraceUID    string `config:"traceUid"`
		TraceReqID  string `config:"traceReqId"`
		TraceSpanID string `config:"traceSpanId"`
		TraceBizID  string `config:"traceBizId"`
	}
)

// DefaultConfig of system
func DefaultConfig() *Config {
	return &Config{
		Name:        "box",
		Version:     "unknown",
		TraceUID:    "box.trace.uid",
		TraceReqID:  "box.trace.reqId",
		TraceSpanID: "box.trace.spanId",
		TraceBizID:  "box.trace.bizId",
	}
}

// StdConfig get from box file
func StdConfig() *Config {
	cfg := DefaultConfig()
	config.Get("box").Scan(cfg)

	return cfg
}

// Build system
func (c *Config) Build() *System {
	return newSystem(c)
}
