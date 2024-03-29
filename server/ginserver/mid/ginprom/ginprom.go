package ginprom

import (
	"strconv"
	"time"

	metric2 "github.com/boxgo/box/v2/metric"
	"github.com/gin-gonic/gin"
)

type (
	GinProm struct {
		cfg                *Config
		reqSizeSummary     *metric2.SummaryVec
		reqBeginCounter    *metric2.CounterVec
		reqFinishCounter   *metric2.CounterVec
		reqDurationSummary *metric2.SummaryVec
		resSizeSummary     *metric2.SummaryVec
	}
)

func newGinProm(c *Config) *GinProm {
	return &GinProm{
		cfg: c,
		reqSizeSummary: metric2.NewSummaryVec(
			"http_server_request_size_bytes",
			"The HTTP request sizes in bytes.",
			[]string{"method", "url"},
			map[float64]float64{
				0.5:  0.05,
				0.75: 0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		),
		reqBeginCounter: metric2.NewCounterVec(
			"http_server_request_begin_total",
			"How many HTTP requests ready to process.",
			[]string{"method", "url"},
		),
		reqFinishCounter: metric2.NewCounterVec(
			"http_server_request_finish_total",
			"How many HTTP requests processed.",
			[]string{"method", "url", "status", "errcode"},
		),
		reqDurationSummary: metric2.NewSummaryVec(
			"http_server_request_duration_seconds",
			"The HTTP request latencies in seconds.",
			[]string{"method", "url", "status", "errcode"},
			map[float64]float64{
				0.5:  0.05,
				0.75: 0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		),
		resSizeSummary: metric2.NewSummaryVec(
			"http_server_response_size_bytes",
			"The HTTP response sizes in bytes.",
			[]string{"method", "url", "status", "errcode"},
			map[float64]float64{
				0.5:  0.05,
				0.75: 0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		),
	}
}

func (prom *GinProm) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		labels := []string{
			ctx.Request.Method,
			prom.cfg.requestURLMappingFn(ctx),
		}

		reqSz := computeApproximateRequestSize(ctx.Request)

		prom.reqSizeSummary.WithLabelValues(labels...).Observe(reqSz)
		prom.reqBeginCounter.WithLabelValues(labels...).Inc()

		ctx.Next()

		resSz := ctx.Writer.Size()
		if resSz < 0 || !ctx.Writer.Written() {
			resSz = 0
		}

		labels = append(labels, strconv.Itoa(ctx.Writer.Status()), strconv.Itoa(ctx.GetInt("errcode")))
		prom.resSizeSummary.WithLabelValues(labels...).Observe(float64(resSz))
		prom.reqFinishCounter.WithLabelValues(labels...).Inc()
		prom.reqDurationSummary.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
	}
}
