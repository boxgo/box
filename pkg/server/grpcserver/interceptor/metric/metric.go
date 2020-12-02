package metric

import (
	"context"
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/errcode"
	"github.com/boxgo/box/pkg/metric"
	"google.golang.org/grpc"
)

var (
	handledCounter = metric.NewCounterVec(
		"grpc_server_handled_total",
		"gGPC server handle msg count",
		[]string{"method", "type", "code"},
	)
	handledSeconds = metric.NewSummaryVec(
		"grpc_server_handled_second",
		"gGPC server handle msg duration",
		[]string{"method", "type", "code"},
		map[float64]float64{
			0.25: 0.05,
			0.5:  0.05,
			0.75: 0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	)
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		report(info.FullMethod, "unary", start, err)

		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		err := handler(srv, ss)

		if info.IsClientStream && info.IsServerStream {
			report(info.FullMethod, "stream_bidi", start, err)
		} else if info.IsClientStream {
			report(info.FullMethod, "stream_client", start, err)
		} else if info.IsServerStream {
			report(info.FullMethod, "stream_server", start, err)
		}

		return err
	}
}

func report(method, typ string, start time.Time, err error) {
	var labels []string

	if err != nil {
		labels = []string{method, typ, fmt.Sprintf("%d", errcode.ParseStatus(err).Code)}
	} else {
		labels = []string{method, typ, "0"}
	}

	handledCounter.WithLabelValues(labels...).Inc()
	handledSeconds.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
}
