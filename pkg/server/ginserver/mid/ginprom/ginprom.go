package ginprom

import (
	"strconv"
	"time"

	"github.com/boxgo/box/pkg/metric"
	"github.com/gin-gonic/gin"
)

type (
	GinProm struct {
		cfg *Config
	}
)

var (
	// 1. Saturation: 饱和度 (Requests Inflight)
	// 衡量服务当前的忙碌程度，通常使用正在处理的请求数来表示。
	reqInFlight = metric.NewGaugeVec(
		"http_server_requests_inflight",
		"The number of HTTP requests currently being processed.",
		[]string{"method", "url"},
	)

	// 2. Traffic & Errors (Counter) -> 负责“高基数”统计
	// 专门用于统计详细的业务错误码 (errcode)。
	// Counter 存储成本低，这里放 errcode 是划算的。
	reqTotal = metric.NewCounterVec(
		"http_server_requests_total",
		"The total number of HTTP requests processed.",
		[]string{"method", "url", "status", "errcode"},
	)

	// 3. Traffic (Size Histogram)
	// 响应大小分布通常只关注正常请求，或者只关注接口维度的带宽压力。
	reqSize = metric.NewHistogramVec(
		"http_server_request_size_bytes",
		"The HTTP request body sizes in bytes.",
		[]string{"method", "url"},
		// 1KB, 5KB, 10KB, 100KB, 1MB, 10MB
		[]float64{1024, 5120, 10240, 102400, 1048576, 10485760},
	)
	resSize = metric.NewHistogramVec(
		"http_server_response_size_bytes",
		"The HTTP response body sizes in bytes.",
		[]string{"method", "url"},
		// 1KB, 5KB, 10KB, 100KB, 1MB, 10MB
		[]float64{1024, 5120, 10240, 102400, 1048576, 10485760},
	)

	// 4. Latency (Duration Histogram) -> 负责“耗时分布”统计
	// 关键修改：移除 errcode。
	// 原因：我们主要通过 status (如 500) 或 url 来定位慢请求。
	// 移除 errcode 后，Series 数量将减少几十倍，极大地降低 Prometheus 压力。
	// 如果真有“特定业务错误码导致慢请求”的极端场景，通常通过 Log/Trace 系统排查，而不是 Metric。
	reqDuration = metric.NewHistogramVec(
		"http_server_request_duration_seconds",
		"The HTTP request latencies in seconds.",
		[]string{"method", "url", "status"},
		[]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	)
)

func newGinProm(c *Config) *GinProm {
	return &GinProm{
		cfg: c,
	}
}

func (prom *GinProm) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		labels := []string{
			ctx.Request.Method,
			prom.cfg.requestURLMappingFn(ctx),
		}

		// Saturation: +1
		reqInFlight.WithLabelValues(labels...).Inc()
		defer reqInFlight.WithLabelValues(labels...).Dec()

		// Traffic: Request Size
		reqSz := computeApproximateRequestSize(ctx.Request)
		reqSize.WithLabelValues(labels...).Observe(reqSz)

		ctx.Next()

		resSz := ctx.Writer.Size()
		if resSz < 0 || !ctx.Writer.Written() {
			resSz = 0
		}

		labels = append(labels, strconv.Itoa(ctx.Writer.Status()), strconv.Itoa(ctx.GetInt("errcode")))

		// Traffic: Response Size & Total Count (implies Errors via labels)
		resSize.WithLabelValues(labels...).Observe(float64(resSz))
		reqTotal.WithLabelValues(labels...).Inc()

		// Latency
		reqDuration.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
	}
}
