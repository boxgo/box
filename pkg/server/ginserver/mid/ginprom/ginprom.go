package ginprom

import (
	"strconv"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/metric"
	"github.com/gin-gonic/gin"
)

type (
	GinProm struct {
		cfg                *Config
		reqCounter         *metric.CounterVec
		reqDurationSummary *metric.SummaryVec
	}
)

func newGinProm(c *Config) *GinProm {
	return &GinProm{
		cfg: c,
		reqCounter: metric.NewCounterVec(
			"http_request_total",
			"How many HTTP requests processed, partitioned by status code and HTTP method.",
			[]string{"status", "retcode", "method", "url", "handler"},
		),
		reqDurationSummary: metric.NewSummaryVec(
			"http_request_duration_seconds",
			"The HTTP request latencies in seconds.",
			[]string{"status", "retcode", "method", "url", "handler"},
			map[float64]float64{
				0.25: 0.05,
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

		ctx.Next()

		labels := []string{
			strconv.Itoa(ctx.Writer.Status()),
			strconv.Itoa(ctx.GetInt("retcode")),
			ctx.Request.Method,
			prom.cfg.requestURLMappingFn(ctx),
			getHandlerName(ctx),
		}

		prom.reqCounter.WithLabelValues(labels...).Inc()
		prom.reqDurationSummary.WithLabelValues(labels...).Observe(time.Since(start).Seconds())
	}
}

func getHandlerName(c *gin.Context) string {
	longName := c.HandlerName()
	shortName := longName

	if idx := strings.LastIndex(longName, "."); idx != -1 {
		shortName = longName[idx+1:]
	}
	return shortName
}
