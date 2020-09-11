package wukong

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boxgo/box/pkg/client/wukong/util"
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
		body       Body
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
		BindBody(&body).
		Error()

	util.AssertEqual(t, err, nil)
	util.AssertEqual(t, statusCode, 200)
	util.AssertEqual(t, statusMsg, "200 OK")
	util.AssertEqual(t, isTimeout, false)
	util.AssertEqual(t, isCancel, false)
}
