package wukong

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boxgo/box/v2/util/testutil"
)

func TestSimple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
	}))
	defer ts.Close()

	client := New("")

	type Body struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
	}

	var (
		statusCode int
		statusMsg  string
		bodyData   Body
		header     http.Header
		isTimeout  bool
		isCancel   bool
	)
	err := client.Get(ts.URL).End().
		BindStatusCode(&statusCode).
		BindStatus(&statusMsg).
		BindHeader(&header).
		BindIsTimeout(&isTimeout).
		BindIsCancel(&isCancel).
		BindBody(&bodyData).
		Error()

	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, statusCode, 200)
	testutil.ExpectEqual(t, statusMsg, "200 OK")
	testutil.ExpectEqual(t, isTimeout, false)
	testutil.ExpectEqual(t, isCancel, false)
}

func TestSample(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
	}))
	defer ts.Close()

	client := New(ts.URL)

	type Body struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
	}

	var (
		statusCode int
		statusMsg  string
		bodyData   Body
		header     http.Header
		isTimeout  bool
		isCancel   bool
	)
	err := client.Get("/").WithCTX(context.Background()).End().
		BindStatusCode(&statusCode).
		BindStatus(&statusMsg).
		BindHeader(&header).
		BindIsTimeout(&isTimeout).
		BindIsCancel(&isCancel).
		BindBody(&bodyData).
		Error()

	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, statusCode, 200)
	testutil.ExpectEqual(t, statusMsg, "200 OK")
	testutil.ExpectEqual(t, isTimeout, false)
	testutil.ExpectEqual(t, isCancel, false)
}

func TestUseBefore(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`123`))
	}))
	defer ts.Close()

	client := New(ts.URL).
		UseBefore(func(req *Request) error {
			return errors.New("before")
		}).
		UseAfter(func(req *Request, resp *Response) error {
			return nil
		})

	var (
		err        error
		statusCode int
		statusMsg  string
		header     http.Header
		isTimeout  bool
		isCancel   bool
	)
	client.Get("/").WithCTX(context.Background()).End().
		BindStatusCode(&statusCode).
		BindStatus(&statusMsg).
		BindHeader(&header).
		BindIsTimeout(&isTimeout).
		BindIsCancel(&isCancel).
		BindError(&err)

	testutil.ExpectEqual(t, err.Error(), "before")
	testutil.ExpectEqual(t, statusCode, 0)
	testutil.ExpectEqual(t, statusMsg, "")
	testutil.ExpectEqual(t, isTimeout, false)
	testutil.ExpectEqual(t, isCancel, false)
}

func TestUseAfter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(`123`))
	}))
	defer ts.Close()

	client := New(ts.URL).
		UseBefore(func(req *Request) error {
			return nil
		}).
		UseAfter(func(req *Request, resp *Response) error {
			return errors.New("after")
		})

	var (
		err        error
		statusCode int
		statusMsg  string
		header     http.Header
		isTimeout  bool
		isCancel   bool
	)
	client.Get("/").WithCTX(context.Background()).End().
		BindStatusCode(&statusCode).
		BindStatus(&statusMsg).
		BindHeader(&header).
		BindIsTimeout(&isTimeout).
		BindIsCancel(&isCancel).
		BindError(&err)

	testutil.ExpectEqual(t, err.Error(), "after")
	testutil.ExpectEqual(t, statusCode, 400)
	testutil.ExpectEqual(t, statusMsg, "400 Bad Request")
	testutil.ExpectEqual(t, isTimeout, false)
	testutil.ExpectEqual(t, isCancel, false)
}
