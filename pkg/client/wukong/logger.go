package wukong

import (
	"github.com/boxgo/box/pkg/logger"
	"moul.io/http2curl"
)

func loggerStart(req *Request) error {
	r, e := req.RawRequest()
	if e != nil {
		return e
	}

	curl, e := http2curl.GetCurlCommand(r)
	if e != nil {
		return e
	}

	logger.Trace(req.Context).Infow("http_request_start", "request", curl.String())

	return nil
}

func loggerAfter(req *Request, resp *Response) error {
	logger.Trace(req.Context).Infow("http_request_end", "response", string(resp.Bytes()))

	return nil
}
