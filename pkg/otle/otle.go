package otle

import (
	"context"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/system"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresour "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitProvider(ctx context.Context, agentAddr string) (func(), error) {
	res, err := sdkresour.New(ctx,
		sdkresour.WithFromEnv(),
		sdkresour.WithProcess(),
		sdkresour.WithTelemetrySDK(),
		sdkresour.WithHost(),
		sdkresour.WithOSType(),
		sdkresour.WithContainerID(),
		sdkresour.WithAttributes(
			semconv.ServiceNamespaceKey.String(config.ServiceNamespace()),
			semconv.ServiceNameKey.String(config.ServiceName()),
			semconv.ServiceVersionKey.String(config.ServiceVersion()),
			semconv.ServiceInstanceIDKey.String(system.Hostname()),
		),
	)

	if err != nil {
		return nil, err
	}

	traceExp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(agentAddr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}, nil
}
