package wukong

import (
	"strconv"
	"time"

	"github.com/boxgo/box/pkg/metric"
	"golang.org/x/net/context"
)

type (
	metricDurationKey struct{}
)

const (
	metricSwitchKey = "metric.enable"
)

var (
	requestInflight = metric.NewGaugeVec(
		"http_client_request_in_process",
		"http client requesting",
		[]string{"method", "baseUrl", "url"},
	)
	requestCounter = metric.NewCounterVec(
		"http_client_request_total",
		"http client request counter",
		[]string{"method", "baseUrl", "url", "statusCode", "error"},
	)
	requestDuration = metric.NewSummaryVec(
		"http_client_request_seconds",
		"http client request duration",
		[]string{"method", "baseUrl", "url", "statusCode", "error"},
		map[float64]float64{
			0.5:  0.05,
			0.75: 0.05,
			0.9:  0.01,
			0.99: 0.001,
			1:    0.001,
		},
	)
)

func metricStart(request *Request) error {
	if val, ok := request.Context.Value(metricSwitchKey).(bool); ok && !val {
		return nil
	}

	requestInflight.WithLabelValues(request.Method, request.BaseUrl, request.Url).Inc()

	request.Context = context.WithValue(request.Context, metricDurationKey{}, time.Now())

	return nil
}

func metricEnd(request *Request, resp *Response) error {
	if val, ok := request.Context.Value(metricSwitchKey).(bool); ok && !val {
		return nil
	}

	var (
		errMsg     = ""
		duration   = time.Duration(0)
		statusCode = strconv.Itoa(resp.StatusCode())
	)

	if resp.Error() != nil {
		errMsg = resp.Error().Error()
	}

	if start, ok := request.Context.Value(metricDurationKey{}).(time.Time); ok {
		duration = time.Since(start)
	}

	requestInflight.WithLabelValues(request.Method, request.BaseUrl, request.Url).Dec()
	requestCounter.WithLabelValues(request.Method, request.BaseUrl, request.Url, statusCode, errMsg).Inc()
	requestDuration.WithLabelValues(request.Method, request.BaseUrl, request.Url, statusCode, errMsg).Observe(duration.Seconds())

	return nil
}
