package zaplog

import (
	"context"
	"strings"

	"github.com/boxgo/box/pkg/component"
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	Logger struct {
		component.NoopBox
		cfg config.SubConfigurator
	}

	Options struct {
		cfg config.SubConfigurator
	}

	OptionFunc func(*Options)
)

var (
	Default = New()
)

func WithConfigurator(cfg config.SubConfigurator) OptionFunc {
	return func(opts *Options) {
		opts.cfg = cfg
	}
}

func New(optionFunc ...OptionFunc) *Logger {
	opts := &Options{}
	for _, fn := range optionFunc {
		fn(opts)
	}

	if opts.cfg == nil {
		opts.cfg = config.Default
	}

	l := &Logger{
		cfg: opts.cfg,
	}

	return l
}

func (l *Logger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)

		newCtx := l.ctx(ctx, md, info.FullMethod)

		logger.TraceRaw(newCtx).Info(">>>", []zap.Field{zap.Any("req", req)}...)

		resp, err = handler(newCtx, req)
		if err != nil {
			logger.TraceRaw(newCtx).Info("xxx", []zap.Field{zap.Any("err", err)}...)
		}

		logger.TraceRaw(newCtx).Info("<<<", []zap.Field{zap.Any("resp", resp)}...)

		return resp, err
	}
}

func (l *Logger) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, _ := metadata.FromIncomingContext(ctx)

		newCtx := l.ctx(ctx, md, info.FullMethod)

		err := handler(newCtx, ss)
		if err != nil {
			logger.TraceRaw(newCtx).Info("xxx", []zap.Field{zap.Any("err", err)}...)
		}

		logger.TraceRaw(newCtx).Info("<<<", []zap.Field{zap.Any("resp", "")}...)

		return err
	}
}

func (l *Logger) ctx(ctx context.Context, md metadata.MD, biz string) context.Context {
	uidKey := l.cfg.GetTraceUid()
	reqIdKey := l.cfg.GetTraceReqId()
	bizIdKey := l.cfg.GetTraceBizId()
	spanKey := l.cfg.GetTraceSpanId()

	uidVal := getFirst(md.Get(strings.ToLower(uidKey)))
	reqIdVal := getFirst(md.Get(strings.ToLower(reqIdKey)))
	spanVal := getFirst(md.Get(strings.ToLower(spanKey)))
	bizIdVal := biz

	newCtx := context.WithValue(ctx, uidKey, uidVal)
	newCtx = context.WithValue(newCtx, reqIdKey, reqIdVal)
	newCtx = context.WithValue(newCtx, bizIdKey, bizIdVal)
	newCtx = context.WithValue(newCtx, spanKey, spanVal)

	return newCtx
}

func getFirst(arr []string) string {
	if len(arr) == 0 {
		return ""
	}

	return arr[0]
}
