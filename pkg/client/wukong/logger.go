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
)

func loggerStart(req *Request) error {
	level := getLoggerLevel(req)

	if level&LoggerDisable != 0 {
		return nil
	}

	if level&LoggerRequest == 0 {
		logger.Trace(req.Context).Infow("http_client_start")
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

	if level&LoggerResponse == 0 {
		logger.Trace(req.Context).Info("http_client_end")
	} else {
		logger.Trace(req.Context).Infow("http_client_end", "elapsed", req.TraceInfo.ElapsedTime, "request", req.curl, "response", string(resp.Bytes()))
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
