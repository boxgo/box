package request

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"reflect"

	"github.com/boxgo/box/pkg/client/wukong/encoder"
	"github.com/boxgo/box/pkg/client/wukong/response"
	"github.com/boxgo/box/pkg/client/wukong/util"
)

type (
	Request struct {
		do          func(*Request) *response.Response
		ctx         context.Context
		error       error
		method      string
		baseUrl     string
		url         string
		contentType string
		header      http.Header
		queryData   url.Values
		formData    url.Values
		paramData   map[string]interface{}
		bodyData    interface{}
		cookies     []*http.Cookie
	}
)

func NewRequest(do func(*Request) *response.Response, method, baseUrl, path string) *Request {
	req := &Request{
		do:          do,
		ctx:         nil,
		error:       nil,
		method:      method,
		baseUrl:     baseUrl,
		url:         path,
		contentType: encoder.MimeTypeJSON,
		header:      http.Header{},
		cookies:     make([]*http.Cookie, 0),
		queryData:   url.Values{},
		formData:    url.Values{},
		paramData:   make(map[string]interface{}),
	}

	return req
}

func (request *Request) WithCTX(ctx context.Context) *Request {
	request.ctx = ctx

	return request
}

// SetHeader set header
func (request *Request) SetHeader(key, value string) *Request {
	request.header.Set(key, value)

	return request
}

func (request *Request) AddCookies(cookies ...*http.Cookie) *Request {
	request.cookies = append(request.cookies, cookies...)

	return request
}

// Param method sets multiple URL path key-value pairs.
//
// For example: http://example.com/users/:uid
//		client.Get("http://example.com/users/:uid").Param(map[string]interface{}{"uid": "123"}).End()
// request target url will be replace to `http://example.com/users/123`
func (request *Request) Param(param map[string]interface{}) *Request {
	request.paramData = param

	return request
}

// Query
// format:
//		1.map[string]interface{} {"key": "value", "key1": 1}
func (request *Request) Query(query interface{}) *Request {
	switch v := reflect.ValueOf(query); v.Kind() {
	case reflect.Map:
		request.queryMapOrStruct(v.Interface())
	case reflect.Struct:
		request.queryMapOrStruct(v.Interface())
	}

	return request
}

func (request *Request) Send(data interface{}) *Request {
	request.bodyData = data

	return request
}

func (request *Request) Type(typ string) *Request {
	request.contentType = typ

	return request
}

func (request *Request) End() *response.Response {
	return request.do(request)
}

func (request *Request) RawRequest() (*http.Request, error) {
	var (
		err    error
		req    *http.Request
		reader io.Reader
	)
	if err = request.error; err != nil {
		return req, err
	}

	targetUrl, err := util.UrlJoin(request.baseUrl, request.url)
	if err != nil {
		return req, err
	}

	if request.bodyData != nil {
		if reader, err = encoder.Decode(request.contentType, request.bodyData); err != nil {
			return req, err
		}
	}

	if req, err = http.NewRequest(request.method, util.UrlFormat(targetUrl, request.paramData), reader); err != nil {
		return req, err
	}

	if request.ctx != nil {
		req = req.WithContext(httptrace.WithClientTrace(request.ctx, trace))
	}

	req.Header = request.header
	req.Header.Set("Content-Type", string(request.contentType))

	for _, cookie := range request.cookies {
		req.AddCookie(cookie)
	}

	if len(request.queryData) != 0 {
		req.URL.RawQuery = request.queryData.Encode()
	}

	return req, err
}

func (request *Request) queryMapOrStruct(query interface{}) {
	if marshalContent, err := json.Marshal(query); err != nil {
		request.error = err
	} else {
		var val map[string]interface{}
		if err := json.Unmarshal(marshalContent, &val); err != nil {
			request.error = err
		} else {
			for k, v := range val {
				switch val := v.(type) {
				case []interface{}:
					for _, e := range val {
						request.queryData.Add(k, fmt.Sprintf("%v", e))
					}
				default:
					request.queryData.Add(k, fmt.Sprintf("%v", v))
				}
			}
		}
	}
}
