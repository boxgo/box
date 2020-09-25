package wukong

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

	AssertEqual(t, err, nil)
	AssertEqual(t, statusCode, 200)
	AssertEqual(t, statusMsg, "200 OK")
	AssertEqual(t, isTimeout, false)
	AssertEqual(t, isCancel, false)
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

	AssertEqual(t, err.Error(), "before")
	AssertEqual(t, statusCode, 0)
	AssertEqual(t, statusMsg, "")
	AssertEqual(t, isTimeout, false)
	AssertEqual(t, isCancel, false)
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

	AssertEqual(t, err.Error(), "after")
	AssertEqual(t, statusCode, 400)
	AssertEqual(t, statusMsg, "400 Bad Request")
	AssertEqual(t, isTimeout, false)
	AssertEqual(t, isCancel, false)
}