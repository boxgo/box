package wukong

import (
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/boxgo/box/pkg/client/wukong/request"
	"github.com/boxgo/box/pkg/client/wukong/response"
)

type (
	WuKong struct {
		baseUrl     string
		timeout     time.Duration
		client      *http.Client
		clientTrace *httptrace.ClientTrace
	}
)

func New(baseUrl string) *WuKong {
	w := &WuKong{
		baseUrl: baseUrl,
		client: &http.Client{
			Transport: DefaultTransport,
		},
	}

	return w
}

func (wk *WuKong) SetTransport(transport *http.Transport) *WuKong {
	wk.client.Transport = transport

	return wk
}

func (wk *WuKong) Get(path string) *request.Request {
	return request.NewRequest(wk.do, http.MethodGet, wk.baseUrl, path)
}

func (wk *WuKong) Post(path string) *request.Request {
	return request.NewRequest(wk.do, http.MethodPost, wk.baseUrl, path)
}

func (wk *WuKong) Put(path string) *request.Request {
	return request.NewRequest(wk.do, http.MethodPut, wk.baseUrl, path)
}

func (wk *WuKong) Patch(path string) *request.Request {
	return request.NewRequest(wk.do, http.MethodPatch, wk.baseUrl, path)
}

func (wk *WuKong) Delete(path string) *request.Request {
	return request.NewRequest(wk.do, http.MethodDelete, wk.baseUrl, path)
}

func (wk *WuKong) Head(path string) *request.Request {
	return request.NewRequest(wk.do, http.MethodHead, wk.baseUrl, path)
}

func (wk *WuKong) Options(path string) *request.Request {
	return request.NewRequest(wk.do, http.MethodOptions, wk.baseUrl, path)
}

func (wk *WuKong) Trace(path string) *request.Request {
	return request.NewRequest(wk.do, http.MethodTrace, wk.baseUrl, path)
}

func (wk *WuKong) Timeout(t time.Duration) *WuKong {
	wk.client.Timeout = t

	return wk
}

func (wk *WuKong) do(r *request.Request) *response.Response {
	req, err := r.RawRequest()

	if err != nil {
		return response.New(err, nil)
	}

	resp, err := wk.client.Do(req)

	return response.New(err, resp)
}
