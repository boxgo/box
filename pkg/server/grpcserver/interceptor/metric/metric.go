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
	// Saturation
	reqInflight = metric.NewGaugeVec(
		"grpc_server_requests_inflight",
		"The number of gRPC requests currently being processed.",
		[]string{"method", "type"},
	)

	// Traffic & Errors
	reqTotal = metric.NewCounterVec(
		"grpc_server_requests_total",
		"The total number of gRPC requests processed.",
		[]string{"method", "type", "code"},
	)

	// Latency
	reqDuration = metric.NewHistogramVec(
		"grpc_server_request_duration_seconds",
		"The gRPC request latencies in seconds.",
		[]string{"method", "type", "code"},
		// .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
		[]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	)
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		typ := "unary"

		reqInflight.WithLabelValues(info.FullMethod, typ).Inc()
		defer reqInflight.WithLabelValues(info.FullMethod, typ).Dec()

		resp, err := handler(ctx, req)

		report(info.FullMethod, typ, start, err)

		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		typ := "stream"
		if info.IsClientStream && info.IsServerStream {
			typ = "stream_bidi"
		} else if info.IsClientStream {
			typ = "stream_client"
		} else if info.IsServerStream {
			typ = "stream_server"
		}

		reqInflight.WithLabelValues(info.FullMethod, typ).Inc()
		defer reqInflight.WithLabelValues(info.FullMethod, typ).Dec()

		err := handler(srv, ss)

		report(info.FullMethod, typ, start, err)

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

	reqTotal.WithLabelValues(labels...).Inc()
	reqDuration.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
}
