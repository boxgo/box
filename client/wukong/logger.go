package wukong

import (
	"github.com/boxgo/box/v2/logger"
	"moul.io/http2curl"
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
)

func loggerStart(req *Request) error {
	level := getLoggerLevel(req)

	if level&LoggerDisable == 1 {
		return nil
	}

	if level&LoggerRequest == 0 {
		logger.Trace(req.Context).Info("http_client_start")
	} else {
		r, e := req.RawRequest()
		if e != nil {
			return e
		}

		curl, e := http2curl.GetCurlCommand(r)
		if e != nil {
			return e
		}

		logger.Trace(req.Context).Infow("http_client_start", "request", curl.String())
	}

	return nil
}

func loggerAfter(req *Request, resp *Response) error {
	level := getLoggerLevel(req)

	if level&LoggerDisable == 1 {
		return nil
	}

	if level&LoggerResponse == 0 {
		logger.Trace(req.Context).Info("http_client_end")
	} else {
		logger.Trace(req.Context).Infow("http_client_end", "response", string(resp.Bytes()))
	}

	return nil
}

func getLoggerLevel(req *Request) LoggerLevel {
	level, ok := req.Context.Value(loggerLevelKey).(LoggerLevel)
	if !ok {
		level = LoggerRequest | LoggerResponse
	}

	return level
}
