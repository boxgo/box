package recovery

import (
	"context"
	"fmt"

	"github.com/boxgo/box/pkg/errcode"
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/metric"
	"google.golang.org/grpc"
)

var (
	panicCounter = metric.NewCounterVec(
		"grpc_server_panic_total",
		"grpc server panic counter",
		[]string{"method", "panic"},
	)
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		panicked := true

		defer func() {
			if panicErr := recover(); panicErr != nil || panicked {
				logger.Errorw("grpc unary server panic:", "panicked", panicked, "panic", panicErr)

				err = errcode.ErrGRPCServerPanic.Build(panicErr)
				panicCounter.WithLabelValues(info.FullMethod, fmt.Sprintf("%s", panicErr)).Inc()
			}
		}()

		resp, err = handler(ctx, req)
		panicked = false

		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true

		defer func() {
			if panicErr := recover(); panicErr != nil || panicked {
				logger.Errorw("grpc stream server panic:", "panicked", panicked, "panic", panicErr)

				err = errcode.ErrGRPCServerPanic.Build(panicErr)
				panicCounter.WithLabelValues(info.FullMethod, fmt.Sprintf("%s", panicErr)).Inc()
			}
		}()

		err = handler(srv, stream)
		panicked = false

		return err
	}
}
