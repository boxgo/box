package wukong

import (
	"context"
	"errors"
	"net"
	"os"
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
		errMsg = classifyHTTPError(resp.Error())
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

// classifyHTTPError 将 HTTP 客户端错误分类为有限的几个类别，避免指标爆炸
// 同时尽可能保留有用的错误信息
// 注意：HTTP 状态码已通过 status 字段上报，此处不再根据状态码分类
func classifyHTTPError(err error) string {
	if err == nil {
		return ""
	}

	// 检查标准库超时错误
	if os.IsTimeout(err) {
		return "timeout_error"
	}

	// 检查 context 超时错误
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
		return "timeout_error"
	}

	// 检查 net.Error 接口的 Timeout() 方法
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return "timeout_error"
		}
		// 如果是网络错误但不是超时，归类为连接错误
		return "connection_error"
	}

	errStr := strings.ToLower(err.Error())

	// 检查 DNS 相关错误
	if isDNSError(errStr) {
		return "dns_error"
	}

	// 检查 TLS/SSL 相关错误
	if isTLSError(errStr) {
		return "tls_error"
	}

	// 检查连接相关错误
	if isHTTPConnectionError(errStr) {
		return "connection_error"
	}

	// 其他错误统一归类
	return "other_error"
}

// isDNSError 判断是否为 DNS 相关错误
func isDNSError(errStr string) bool {
	dnsKeywords := []string{
		"no such host",
		"no hosts found",
		"dns",
		"lookup",
		"unknown host",
		"host not found",
		"name resolution",
		"getaddrinfo",
	}

	for _, keyword := range dnsKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isTLSError 判断是否为 TLS/SSL 相关错误
func isTLSError(errStr string) bool {
	tlsKeywords := []string{
		"tls",
		"ssl",
		"certificate",
		"x509",
		"handshake failure",
		"bad certificate",
		"certificate verify failed",
		"unknown authority",
		"certificate signed by unknown authority",
		"tls:",
		"remote error",
	}

	for _, keyword := range tlsKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isHTTPConnectionError 判断是否为连接相关错误
func isHTTPConnectionError(errStr string) bool {
	connectionKeywords := []string{
		"connection",
		"connect",
		"connection refused",
		"connection reset",
		"connection lost",
		"connection closed",
		"no connection",
		"broken pipe",
		"network",
		"dial tcp",
		"connection timeout",
		"i/o error",
		"use of closed network connection",
		"connection reset by peer",
		"no route to host",
		"refused",
		"closed",
		"EOF",
		"unreachable",
		"network is unreachable",
	}

	for _, keyword := range connectionKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

