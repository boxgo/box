package wukong

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
)

type (
	Response struct {
		err             error
		req             *Request
		resp            *http.Response
		bodyData        []byte
		ignoreEncodeErr bool
	}
)

func NewResponse(err error, req *Request, resp *http.Response) *Response {
	if err != nil {
		return &Response{
			err:  err,
			req:  req,
			resp: resp,
		}
	}

	b := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(resp.Body, b)

	body, err := ioutil.ReadAll(reader)

	defer resp.Body.Close()

	resp.Body = ioutil.NopCloser(b)

	return &Response{
		err:      err,
		req:      req,
		resp:     resp,
		bodyData: body,
	}
}

func (resp *Response) IgnoreEncodeErr() *Response {
	resp.ignoreEncodeErr = true

	return resp
}

func (resp *Response) Error() error {
	return resp.err
}

func (resp *Response) BindError(err *error) *Response {
	if err != nil {
		*err = resp.Error()
	}

	return resp
}

func (resp *Response) Bytes() []byte {
	return resp.bodyData
}

func (resp *Response) BindBytes(b *[]byte) *Response {
	if b != nil {
		*b = resp.bodyData
	}

	return resp
}

func (resp *Response) StatusCode() int {
	if resp.resp == nil {
		return 0
	}

	return resp.resp.StatusCode
}

func (resp *Response) BindStatusCode(code *int) *Response {
	if code != nil {
		*code = resp.StatusCode()
	}

	return resp
}

func (resp *Response) Status() string {
	if resp.resp == nil {
		return ""
	}

	return resp.resp.Status
}

func (resp *Response) BindStatus(status *string) *Response {
	if status != nil {
		*status = resp.Status()
	}

	return resp
}

func (resp *Response) Header() http.Header {
	if resp.resp == nil {
		return http.Header{}
	}

	return resp.resp.Header
}

func (resp *Response) BindHeader(header *http.Header) *Response {
	if header != nil {
		*header = resp.Header()
	}

	return resp
}

func (resp *Response) IsTimeout() bool {
	if resp.err == nil {
		return false
	}

	if strings.Contains(resp.err.Error(), "context deadline exceeded") ||
		strings.Contains(resp.err.Error(), "net/http: timeout") {
		return true
	}

	return false
}

func (resp *Response) BindIsTimeout(ok *bool) *Response {
	if ok != nil {
		*ok = resp.IsTimeout()
	}

	return resp
}

func (resp *Response) IsCancel() bool {
	if resp.err == nil {
		return false
	}

	if strings.Contains(resp.err.Error(), "context canceled") {
		return true
	}

	return false
}

func (resp *Response) BindIsCancel(ok *bool) *Response {
	if ok != nil {
		*ok = resp.IsCancel()
	}

	return resp
}

func (resp *Response) Body() io.ReadCloser {
	if resp.resp == nil {
		return nil
	}

	return resp.resp.Body
}

func (resp *Response) BindBody(data interface{}) *Response {
	if resp.err != nil || data == nil {
		return resp
	}

	contentType := resp.contentType()
	if err := Encode(contentType, resp.bodyData, data); err != nil && !resp.ignoreEncodeErr {
		resp.err = err
	}

	return resp
}

func (resp *Response) ConditionBindBody(check func(interface{}) bool, data ...interface{}) *Response {
	if resp.err != nil || len(data) == 0 {
		return resp
	}

	contentType := resp.contentType()
	for _, d := range data {
		if err := Encode(contentType, resp.bodyData, d); err != nil && !resp.ignoreEncodeErr {
			resp.err = err
			break
		}

		if check(d) {
			break
		}
	}

	return resp
}

func (resp *Response) contentType() string {
	contentType := resp.resp.Header.Get("Content-Type")

	if mediaType, _, err := mime.ParseMediaType(contentType); err != nil || mediaType == "" {
		return resp.req.ContentType
	} else {
		return mediaType
	}
}
