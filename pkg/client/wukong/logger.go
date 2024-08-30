package wukong

import (
	"github.com/boxgo/box/pkg/logger"
)

type (
	LoggerLevel int
)

const (
	loggerLevelKey = "logger.level"
)

const (
	LoggerDisable LoggerLevel = 1 << iota
	LoggerRequest
	LoggerResponse
	LoggerCurl
	LoggerTrace
)

func loggerStart(req *Request) error {
	level := getLoggerLevel(req)

	if level&LoggerDisable != 0 {
		return nil
	}

	if level&LoggerRequest == 0 {
		logger.Trace(req.Context).Infow("http_client_start", "baseUrl", req.BaseUrl, "url", req.Url)
	} else {
		logger.Trace(req.Context).Infow("http_client_start", "request", req.curl)
	}

	return nil
}

func loggerAfter(req *Request, resp *Response) error {
	level := getLoggerLevel(req)

	if level&LoggerDisable != 0 {
		return nil
	}

	if level&LoggerTrace != 0 {
		logger.Trace(req.Context).Infow("http_client_trace", "baseUrl", req.BaseUrl, "url", req.Url,
			"HostPort", req.TraceInfo.HostPort,
			"ConnectNetwork", req.TraceInfo.ConnectNetwork,
			"ConnectAddr", req.TraceInfo.ConnectAddr,
			"ElapsedTime", req.TraceInfo.ElapsedTime,
			"GetConnElapsed", req.TraceInfo.GetConnElapsed,
			"DNSLookupElapsed", req.TraceInfo.DNSLookupElapsed,
			"ConnectElapsed", req.TraceInfo.ConnectElapsed,
			"TLSHandshakeElapsed", req.TraceInfo.TLSHandshakeElapsed,
		)
	}

	if level&LoggerResponse == 0 {
		logger.Trace(req.Context).Infow("http_client_end", "baseUrl", req.BaseUrl, "url", req.Url)
	} else {
		logger.Trace(req.Context).Infow("http_client_end", "elapsed", req.TraceInfo.ElapsedTime.Milliseconds(), "request", req.curl, "response", string(resp.Bytes()))
	}

	return nil
}

func getLoggerLevel(req *Request) LoggerLevel {
	level, ok := req.Context.Value(loggerLevelKey).(LoggerLevel)
	if !ok {
		level = req.client.logger
	}

	return level
}
