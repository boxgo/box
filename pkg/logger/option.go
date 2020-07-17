package logger

import (
	"github.com/boxgo/box/pkg/config"
)

type (
	Options struct {
		Level          *config.Field
		Encoding       *config.Field
		TraceUid       *config.Field
		TraceRequestId *config.Field
		TraceSpanId    *config.Field
		TraceBizId     *config.Field
	}
)

func NewOptions(name string) *Options {
	return &Options{
		Level:          config.NewField(name, "level", "Levels: debug,info,warn,error,dpanic,panic,fatal", "info"),
		Encoding:       config.NewField(name, "encoding", "PS: console or json", "console"),
		TraceUid:       config.NewField(name, "traceUid", "Name as trace uid in context", "uid"),
		TraceRequestId: config.NewField(name, "traceRequestId", "Name as trace requestId in context", "requestId"),
		TraceSpanId:    config.NewField(name, "traceSpanId", "Name as trace spanId in context", "traceSpanId"),
		TraceBizId:     config.NewField(name, "traceBizId", "Name as trace bizId in context", "traceBizId"),
	}
}

func (opts *Options) Fields() []*config.Field {
	return []*config.Field{
		opts.Level,
		opts.Encoding,
		opts.TraceUid,
		opts.TraceRequestId,
		opts.TraceSpanId,
		opts.TraceBizId,
	}
}
