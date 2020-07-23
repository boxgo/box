package reqid

import (
	"context"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/server/grpcserver/interceptor/wrapper"
	"github.com/boxgo/box/pkg/util/strutil"
	"github.com/teris-io/shortid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	RequestId struct {
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

func New(optFunc ...OptionFunc) *RequestId {
	opts := &Options{}
	for _, fn := range optFunc {
		fn(opts)
	}

	if opts.cfg == nil {
		opts.cfg = config.Default
	}

	return &RequestId{
		cfg: opts.cfg,
	}
}

func (reqId *RequestId) UnaryRequestId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		newCtx := reqId.newCtx(ctx, &md)

		return handler(newCtx, req)
	}
}

func (reqId *RequestId) StreamRequestId() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(srv, ss)
		}

		wrapped := wrapper.WrapServerStream(ss)
		wrapped.WrappedContext = reqId.newCtx(ctx, &md)

		return handler(srv, wrapped)
	}
}

func (reqId *RequestId) newCtx(ctx context.Context, md *metadata.MD) context.Context {
	return context.WithValue(ctx, reqId.cfg.GetTraceReqId(), reqId.getReqId(md))
}

func (reqId *RequestId) getReqId(md *metadata.MD) (str string) {
	if reqIdArr := md.Get(reqId.cfg.GetTraceReqId()); len(reqIdArr) > 0 && reqIdArr[0] != "" {
		str = reqIdArr[0]
	} else if id, err := shortid.Generate(); err == nil {
		str = id
	} else {
		str = strutil.RandString(10)
	}

	return
}
