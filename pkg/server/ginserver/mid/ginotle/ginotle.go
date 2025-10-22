package ginotle

import (
	"fmt"

	"github.com/boxgo/box/pkg/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

type (
	GinOTLE struct {
		cfg               *Config
		tracer            trace.Tracer
		textMapPropagator propagation.TextMapPropagator
	}
)

func newGinOTLE(c *Config) *GinOTLE {
	tracer := otel.Tracer("")
	propagators := otel.GetTextMapPropagator()

	return &GinOTLE{
		cfg:               c,
		tracer:            tracer,
		textMapPropagator: propagators,
	}
}

func (prom *GinOTLE) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		ctx := prom.textMapPropagator.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []trace.SpanStartOption{
			trace.WithAttributes(httpconv.ServerRequest(config.ServiceName(), c.Request)...),
			trace.WithSpanKind(trace.SpanKindServer),
		}
		var spanName string
		spanName = c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		} else {
			rAttr := semconv.HTTPRoute(spanName)
			opts = append(opts, trace.WithAttributes(rAttr))
		}

		traceID := c.Request.Header.Get("trace-id")
		spanID := c.Request.Header.Get("span-id")

		traceId, _ := trace.TraceIDFromHex(traceID)
		spanId, _ := trace.SpanIDFromHex(spanID)

		spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceId,
			SpanID:     spanId,
			TraceFlags: trace.FlagsSampled, // 这个没写，是不会记录的
			TraceState: trace.TraceState{},
			Remote:     true,
		})
		// 不用 pctx，不会把 spanctx 当做 parentCtx
		ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)

		ctx, span := prom.tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()

		status := c.Writer.Status()
		span.SetStatus(httpconv.ServerStatus(status))
		if status > 0 {
			span.SetAttributes(semconv.HTTPStatusCode(status))
		}
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
	}
}
