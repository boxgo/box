package wukong

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/metric"
)

type (
	metricDurationKey struct{}
)

const (
	metricSwitchKey = "metric.enable"
)

var (
	requestInflight = metric.NewGaugeVec(
		"http_client_requests_inflight",
		"The number of HTTP client requests currently in flight.",
		[]string{"method", "baseUrl", "url"},
	)
	requestCounter = metric.NewCounterVec(
		"http_client_requests_total",
		"The total number of HTTP client requests sent.",
		[]string{"method", "baseUrl", "url", "status", "error"},
	)
	requestDuration = metric.NewHistogramVec(
		"http_client_request_duration_seconds",
		"The HTTP client request latencies in seconds.",
		[]string{"method", "baseUrl", "url", "status", "error"},
		// 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s
		[]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	)
)

// stripQuery removes query parameters and fragment from URL
func stripQuery(url string) string {
	if idx := strings.IndexAny(url, "?#"); idx != -1 {
		return url[:idx]
	}
	return url
}

func metricStart(request *Request) error {
	if val, ok := request.Context.Value(metricSwitchKey).(bool); ok && !val {
		return nil
	}

	url := stripQuery(request.Url)
	requestInflight.WithLabelValues(request.Method, request.BaseUrl, url).Inc()
	request.Context = context.WithValue(request.Context, metricDurationKey{}, time.Now())

	return nil
}

func metricEnd(request *Request, resp *Response) error {
	if val, ok := request.Context.Value(metricSwitchKey).(bool); ok && !val {
		return nil
	}

	var (
		errMsg   = ""
		duration = time.Duration(0)
		status   = strconv.Itoa(resp.StatusCode())
	)

	if resp.Error() != nil {
		errMsg = "error"
	}

	if start, ok := request.Context.Value(metricDurationKey{}).(time.Time); ok {
		duration = time.Since(start)
	}

	url := stripQuery(request.Url)
	requestInflight.WithLabelValues(request.Method, request.BaseUrl, url).Dec()
	requestCounter.WithLabelValues(request.Method, request.BaseUrl, url, status, errMsg).Inc()
	requestDuration.WithLabelValues(request.Method, request.BaseUrl, url, status, errMsg).Observe(duration.Seconds())

	return nil
}
