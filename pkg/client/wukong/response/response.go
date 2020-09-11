package response

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/boxgo/box/pkg/client/wukong/encoder"
)

type (
	Response struct {
		err  error
		resp *http.Response
	}

	ConditionBind struct {
		Data  interface{}
		Check func() bool
	}
)

func New(err error, resp *http.Response) *Response {
	return &Response{
		err:  err,
		resp: resp,
	}
}

func (resp *Response) Error() error {
	return resp.err
}

func (resp *Response) BindError(err *error) *Response {
	*err = resp.Error()

	return resp
}

func (resp *Response) StatusCode() int {
	if resp.resp == nil {
		return 0
	}

	return resp.resp.StatusCode
}

func (resp *Response) BindStatusCode(code *int) *Response {
	*code = resp.StatusCode()

	return resp
}

func (resp *Response) Status() string {
	if resp.resp == nil {
		return ""
	}

	return resp.resp.Status
}

func (resp *Response) BindStatus(status *string) *Response {
	*status = resp.Status()

	return resp
}

func (resp *Response) Header() http.Header {
	if resp.resp == nil {
		return http.Header{}
	}

	return resp.resp.Header
}

func (resp *Response) BindHeader(header *http.Header) *Response {
	if resp.resp == nil {
		return resp
	}

	*header = resp.Header()

	return resp
}

func (resp *Response) BindIsTimeout(ok *bool) *Response {
	if resp.err == nil {
		*ok = false
		return resp
	}

	*ok = resp.IsTimeout()

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

func (resp *Response) BindIsCancel(ok *bool) *Response {
	if resp.err == nil {
		*ok = false
		return resp
	}

	*ok = resp.IsCancel()

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

func (resp *Response) Body() io.ReadCloser {
	if resp.resp == nil {
		return nil
	}

	return resp.resp.Body
}

func (resp *Response) BindBody(data interface{}) *Response {
	if resp.err != nil {
		return resp
	}

	contentType := resp.resp.Header.Get("Content-Type")
	if err := encoder.Encode(contentType, resp.resp.Body, data); err != nil {
		resp.err = err
	}

	return resp
}

func (resp *Response) ConditionBindBody(check func(interface{}) bool, data ...interface{}) *Response {
	if resp.err != nil {
		return resp
	}

	contentType := resp.resp.Header.Get("Content-Type")

	body, err := ioutil.ReadAll(resp.resp.Body)
	if err != nil {
		resp.err = err
		return resp
	}

	for _, d := range data {
		if err := encoder.Encode(contentType, bytes.NewBuffer(body), d); err != nil {
			resp.err = err
			break
		}

		if check(d) {
			break
		}
	}

	return resp
}
