package wukong

import (
	"net/http"
	"time"
)

type (
	WuKong struct {
		baseUrl   string
		client    *http.Client
		basicAuth BasicAuth
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
		client: &http.Client{
			Transport: DefaultTransport,
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

func (wk *WuKong) initRequest(request *Request) *Request {
	request.SetBasicAuth(wk.basicAuth)

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
			break
		}
	}

	for i := 0; i < 1; i++ {
		if err != nil {
			resp = NewResponse(err, nil)
			break
		}

		rawReq, err = req.RawRequest()
		if err != nil {
			resp = NewResponse(err, nil)
			break
		}

		rawResp, err = wk.client.Do(rawReq)
		req.TraceInfo.ElapsedTime = time.Since(startAt)

		resp = NewResponse(err, rawResp)
	}

	for _, after := range wk.after {
		if err = after(req, resp); err != nil {
			resp = NewResponse(err, rawResp)
			break
		}
	}

	return resp
}
