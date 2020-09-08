package wukong

import (
	"log"
	"net/http"
	"time"

	"moul.io/http2curl"
)

type (
	WuKong struct {
		baseUrl string
		timeout time.Duration
		client  *http.Client
	}

	Do func(*Request) (*Response, error)
)

func New(baseUrl string) *WuKong {
	w := &WuKong{
		baseUrl: baseUrl,
		client:  &http.Client{},
	}

	return w
}

func (wk *WuKong) Get(path string) *Request {
	return NewRequest(wk.do, http.MethodGet, wk.baseUrl, path)
}

func (wk *WuKong) Post(path string) *Request {
	return NewRequest(wk.do, http.MethodPost, wk.baseUrl, path)
}

func (wk *WuKong) Put(path string) *Request {
	return NewRequest(wk.do, http.MethodPut, wk.baseUrl, path)
}

func (wk *WuKong) Patch(path string) *Request {
	return NewRequest(wk.do, http.MethodPatch, wk.baseUrl, path)
}

func (wk *WuKong) Delete(path string) *Request {
	return NewRequest(wk.do, http.MethodDelete, wk.baseUrl, path)
}

func (wk *WuKong) Head(path string) *Request {
	return NewRequest(wk.do, http.MethodHead, wk.baseUrl, path)
}

func (wk *WuKong) Options(path string) *Request {
	return NewRequest(wk.do, http.MethodOptions, wk.baseUrl, path)
}

func (wk *WuKong) Trace(path string) *Request {
	return NewRequest(wk.do, http.MethodTrace, wk.baseUrl, path)
}

func (wk *WuKong) Timeout(t time.Duration) *WuKong {
	wk.client.Timeout = t

	return wk
}

func (wk *WuKong) do(req *Request) (*Response, error) {
	r, err := req.RawRequest()
	if err != nil {
		return nil, err
	}

	log.Println(http2curl.GetCurlCommand(r))

	resp, err := wk.client.Do(r)
	if err != nil {
		return nil, err
	}

	return &Response{resp: resp}, nil
}
