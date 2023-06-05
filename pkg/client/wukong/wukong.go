package wukong

import (
	"net/http"
	"time"

	"github.com/boxgo/box/pkg/logger"
)

type (
	WuKong struct {
		baseUrl   string
		client    *http.Client
		logger    LoggerLevel
		metric    bool
		basicAuth BasicAuth
		query     map[string]string
		header    map[string]string
		before    []Before
		after     []After
	}

	BasicAuth struct {
		Username string
		Password string
	}

	Before func(*Request) error
	After  func(*Request, *Response) error
)

func New(baseUrl string) *WuKong {
	w := &WuKong{
		baseUrl: baseUrl,
		logger:  LoggerResponse | LoggerRequest | LoggerCurl,
		metric:  true,
		before:  []Before{loggerStart, metricStart},
		after:   []After{loggerAfter, metricEnd},
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 10,
		},
	}

	return w
}

func (wk *WuKong) UseBefore(fns ...Before) *WuKong {
	wk.before = append(wk.before, fns...)

	return wk
}

func (wk *WuKong) UseAfter(fns ...After) *WuKong {
	wk.after = append(wk.after, fns...)

	return wk
}

func (wk *WuKong) SetTransport(transport *http.Transport) *WuKong {
	wk.client.Transport = transport

	return wk
}

func (wk *WuKong) SetBasicAuth(auth BasicAuth) *WuKong {
	wk.basicAuth = auth

	return wk
}

func (wk *WuKong) SetQuery(query map[string]string) *WuKong {
	wk.query = query

	return wk
}

func (wk *WuKong) SetHeader(header map[string]string) *WuKong {
	wk.header = header

	return wk
}

func (wk *WuKong) Logger(lv LoggerLevel) *WuKong {
	wk.logger = lv

	return wk
}

func (wk *WuKong) Metric(enable bool) *WuKong {
	wk.metric = enable

	return wk
}

func (wk *WuKong) Get(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodGet, path))
}

func (wk *WuKong) Post(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodPost, path))
}

func (wk *WuKong) Put(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodPut, path))
}

func (wk *WuKong) Patch(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodPatch, path))
}

func (wk *WuKong) Delete(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodDelete, path))
}

func (wk *WuKong) Head(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodHead, path))
}

func (wk *WuKong) Options(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodOptions, path))
}

func (wk *WuKong) Trace(path string) *Request {
	return wk.initRequest(NewRequest(wk, http.MethodTrace, path))
}

func (wk *WuKong) Timeout(t time.Duration) *WuKong {
	wk.client.Timeout = t

	return wk
}

func (wk *WuKong) Client() *http.Client {
	return wk.client
}

func (wk *WuKong) initRequest(request *Request) *Request {
	request.SetBasicAuth(wk.basicAuth)
	request.Query(wk.query)
	request.Logger(wk.logger)
	request.Metric(wk.metric)

	for k, v := range wk.header {
		request.SetHeader(k, v)
	}

	return request
}

func (wk *WuKong) do(req *Request) (resp *Response) {
	var (
		err     error
		rawReq  *http.Request
		rawResp *http.Response
		startAt = time.Now()
	)

	for _, before := range wk.before {
		if err = before(req); err != nil {
			logger.Trace(req.Context).Errorw("http_before_hook_error", "err", err)
			break
		}
	}

	for i := 0; i < 1; i++ {
		if err != nil {
			resp = NewResponse(err, req, nil)
			break
		}

		rawReq, err = req.RawRequest()
		if err != nil {
			resp = NewResponse(err, req, nil)
			break
		}

		rawResp, err = wk.client.Do(rawReq)
		req.TraceInfo.ElapsedTime = time.Since(startAt)

		resp = NewResponse(err, req, rawResp)
	}

	for _, after := range wk.after {
		if err = after(req, resp); err != nil {
			logger.Trace(req.Context).Errorw("http_after_hook_error", "err", err)
			resp = NewResponse(err, req, rawResp)
			break
		}
	}

	if err != nil {
		logger.Trace(req.Context).Errorw("http_error", "err", err)
	}

	return resp
}
