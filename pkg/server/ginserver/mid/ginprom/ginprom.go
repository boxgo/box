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
	// Saturation: 饱和度 (Requests Inflight)
	// 衡量服务当前的忙碌程度，通常使用正在处理的请求数来表示。
	reqInFlight = metric.NewGaugeVec(
		"http_server_requests_inflight",
		"The number of HTTP requests currently being processed.",
		[]string{"method", "url"},
	)

	// Traffic: 流量 (Request Rate & Size)
	// 衡量服务的吞吐量，通常使用每秒请求数 (QPS) 或带宽 (IOPS) 来表示。
	// 这里包含了请求总数(reqTotal)、请求包大小(reqSize)和响应包大小(resSize)。
	// Errors: 错误 (Error Rate)
	// 衡量请求失败的比例。
	// 通过 reqTotal 指标中的 status 和 errcode 标签来计算错误率。
	reqTotal = metric.NewCounterVec(
		"http_server_requests_total",
		"The total number of HTTP requests processed.",
		[]string{"method", "url", "status", "errcode"},
	)
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
		[]string{"method", "url", "status", "errcode"},
		// 1KB, 5KB, 10KB, 100KB, 1MB, 10MB
		[]float64{1024, 5120, 10240, 102400, 1048576, 10485760},
	)

	// Latency: 延迟 (Request Duration)
	// 衡量服务处理请求所需的时间。
	reqDuration = metric.NewHistogramVec(
		"http_server_request_duration_seconds",
		"The HTTP request latencies in seconds.",
		[]string{"method", "url", "status", "errcode"},
		// 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s
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
