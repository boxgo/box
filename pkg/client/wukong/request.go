package wukong

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/boxgo/box/pkg/codec/json"
	"github.com/boxgo/box/pkg/util/urlutil"
)

type (
	Request struct {
		client      *WuKong
		rawReq      *http.Request
		TraceInfo   TraceInfo
		Context     context.Context
		BasicAuth   BasicAuth
		Error       error
		Method      string
		BaseUrl     string
		Url         string
		ContentType string
		Header      http.Header
		QueryData   url.Values
		FormData    url.Values
		ParamData   map[string]interface{}
		BodyData    interface{}
		Cookies     []*http.Cookie
		Multipart   []MultipartForm
	}

	MultipartForm struct {
		IsFile    bool
		FieldName string
		Filename  string
		Value     io.Reader
	}
)

var (
	jsonMarshaler = json.NewMarshaler()
)

func NewRequest(client *WuKong, method, path string) *Request {
	req := &Request{
		client:      client,
		rawReq:      nil,
		Context:     context.TODO(),
		Error:       nil,
		Method:      method,
		BaseUrl:     client.baseUrl,
		Url:         path,
		ContentType: MimeTypeJSON,
		Header:      http.Header{},
		Cookies:     make([]*http.Cookie, 0),
		QueryData:   url.Values{},
		FormData:    url.Values{},
		ParamData:   make(map[string]interface{}),
	}

	return req
}

func (request *Request) WithCTX(ctx context.Context) *Request {
	request.Context = ctx

	return request
}

func (request *Request) Set(key, value interface{}) *Request {
	request.Context = context.WithValue(request.Context, key, value)

	return request
}

func (request *Request) Logger(level LoggerLevel) *Request {
	return request.Set(loggerLevelKey, level)
}

func (request *Request) Metric(enable bool) *Request {
	return request.Set(metricSwitchKey, enable)
}

func (request *Request) SetBasicAuth(auth BasicAuth) *Request {
	request.BasicAuth = auth

	return request
}

// SetHeader set Header
func (request *Request) SetHeader(key, value string) *Request {
	request.Header.Set(key, value)

	return request
}

func (request *Request) AddCookies(cookies ...*http.Cookie) *Request {
	request.Cookies = append(request.Cookies, cookies...)

	return request
}

// Param Method sets multiple URL path key-value pairs.
//
// For example: http://example.com/users/:uid
//		client.Get("http://example.com/users/:uid").Param(map[string]interface{}{"uid": "123"}).End()
// request target Url will be replace to `http://example.com/users/123`
func (request *Request) Param(param map[string]interface{}) *Request {
	request.ParamData = param

	return request
}

// Query
// format:
//		1.map[string]interface{} {"key": "value", "key1": 1}
func (request *Request) Query(query interface{}) *Request {
	switch v := reflect.ValueOf(query); v.Kind() {
	case reflect.Map, reflect.Struct:
		request.queryMapOrStruct(request.QueryData, v.Interface())
	case reflect.Ptr:
		switch v.Elem().Kind() {
		case reflect.Map, reflect.Struct:
			request.queryMapOrStruct(request.QueryData, v.Interface())
		}
	}

	return request
}

func (request *Request) Form(form interface{}) *Request {
	switch v := reflect.ValueOf(form); v.Kind() {
	case reflect.Map, reflect.Struct:
		request.queryMapOrStruct(request.FormData, v.Interface())
	case reflect.Ptr:
		switch v.Elem().Kind() {
		case reflect.Map, reflect.Struct:
			request.queryMapOrStruct(request.FormData, v.Interface())
		}
	}

	request.Type(MimeTypeFormData)

	return request
}

func (request *Request) Send(data interface{}) *Request {
	request.BodyData = data

	return request
}

func (request *Request) SendFile(fieldName, filename string) *Request {
	request.Multipart = append(request.Multipart, MultipartForm{
		IsFile:    true,
		FieldName: fieldName,
		Filename:  filename,
		Value:     nil,
	})

	return request
}

func (request *Request) SendFileReader(fieldName, filename string, reader io.Reader) *Request {
	request.Multipart = append(request.Multipart, MultipartForm{
		IsFile:    true,
		FieldName: fieldName,
		Filename:  filename,
		Value:     reader,
	})

	return request
}

func (request *Request) Type(typ string) *Request {
	request.ContentType = typ

	return request
}

func (request *Request) End() *Response {
	return request.client.do(request)
}

func (request *Request) RawRequest() (*http.Request, error) {
	var (
		err       error
		targetUrl string
		req       *http.Request
		reader    io.Reader
	)
	if err = request.Error; err != nil {
		return req, err
	}
	if request.rawReq != nil {
		return request.rawReq, nil
	}

	if strings.HasPrefix(request.BaseUrl, "http") {
		if target, err := urlutil.UrlJoin(request.BaseUrl, request.Url); err != nil {
			return req, err
		} else {
			targetUrl = target
		}
	} else {
		targetUrl = request.Url
	}

	if request.BodyData != nil {
		if data, err := Decode(request.ContentType, request.BodyData); err != nil {
			return req, err
		} else {
			reader = bytes.NewReader(data)
		}
	} else if request.Multipart != nil {
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)

		for _, part := range request.Multipart {
			if part.IsFile {
				var (
					file       *os.File
					partWriter io.Writer
				)

				if part.Value == nil {
					if file, err = os.Open(part.Filename); err != nil {
						break
					}
					defer file.Close()

					if partWriter, err = writer.CreateFormFile(part.FieldName, filepath.Base(part.Filename)); err != nil {
						break
					}

					if _, err = io.Copy(partWriter, file); err != nil {
						break
					}
				} else {
					if partWriter, err = writer.CreateFormFile(part.FieldName, filepath.Base(part.Filename)); err != nil {
						break
					}

					if _, err = io.Copy(partWriter, part.Value); err != nil {
						break
					}
				}

				reader = payload
			} else {
				var (
					partWriter io.Writer
				)

				if partWriter, err = writer.CreateFormField(part.FieldName); err != nil {
					break
				}

				if _, err = io.Copy(partWriter, part.Value); err != nil {
					break
				}
			}
		}

		writer.Close()

		if err != nil {
			request.Error = err

			return req, err
		}

		request.ContentType = writer.FormDataContentType()
	}

	if req, err = http.NewRequest(request.Method, urlutil.UrlFormat(targetUrl, request.ParamData), reader); err != nil {
		return req, err
	}

	if request.Context != nil {
		req = req.WithContext(httptrace.WithClientTrace(request.Context, traceGenerator(request)))
	}

	req.Header = request.Header
	req.Header.Set("Content-Type", string(request.ContentType))

	for _, cookie := range request.Cookies {
		req.AddCookie(cookie)
	}

	if len(request.QueryData) != 0 {
		req.URL.RawQuery = request.QueryData.Encode()
	}
	if len(request.FormData) != 0 {
		if req.URL.RawQuery != "" {
			req.URL.RawQuery += "&" + request.FormData.Encode()
		} else {
			req.URL.RawQuery = request.FormData.Encode()
		}
	}

	if request.BasicAuth.Username != "" || request.BasicAuth.Password != "" {
		req.SetBasicAuth(request.BasicAuth.Username, request.BasicAuth.Password)
	}

	request.rawReq = req

	return req, err
}

func (request *Request) queryMapOrStruct(urlVal url.Values, query interface{}) {
	if marshalContent, err := jsonMarshaler.Marshal(query); err != nil {
		request.Error = err
	} else {
		var val map[string]interface{}
		if err := jsonMarshaler.Unmarshal(marshalContent, &val); err != nil {
			request.Error = err
		} else {
			for k, v := range val {
				switch val := v.(type) {
				case []interface{}:
					for _, e := range val {
						urlVal.Add(k, fmt.Sprintf("%v", e))
					}
				case nil:
					continue
				case string:
					urlVal.Add(k, val)
				case float64:
					urlVal.Add(k, strconv.FormatFloat(val, 'f', -1, 64))
				default:
					if j, err := jsonMarshaler.Marshal(v); err != nil {
						continue
					} else {
						urlVal.Add(k, string(j))
					}
				}
			}
		}
	}
}
