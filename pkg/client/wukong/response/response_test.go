package response

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/client/wukong/util"
)

func TestResponseIsTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	resp := do(req)

	util.AssertEqual(t, true, resp.IsTimeout())
}

func TestResponseIsCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 500)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		cancel()
	}()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	resp := do(req)

	util.AssertEqual(t, true, resp.IsCancel())
}

func TestResponseBindJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type data struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
	}
	d := &data{}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	resp := do(req)
	resp.BindBody(d)

	util.AssertEqual(t, resp.Error(), nil)
	util.AssertEqual(t, d.String, "string")
	util.AssertEqual(t, d.Int, 1)
	util.AssertEqual(t, d.Float, 2.3)
	util.AssertEqual(t, d.Bool, true)
}

func TestResponseConditionBindJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
		} else {
			w.Write([]byte(`{"errcode":1,"errmsg": "not ok"}`))
		}
	}))
	defer ts.Close()

	type Err struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	type Data struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
	}

	condition := func(d interface{}) (ok bool) {
		switch v := d.(type) {
		case *Err:
			ok = v.ErrCode != 0 && v.ErrMsg != ""
		case *Data:
			ok = true
		default:
			ok = false
		}

		return ok
	}

	{
		var (
			err  Err
			data Data
		)
		req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
		resp := do(req)
		e := resp.ConditionBindBody(condition, &err, &data).Error()

		util.AssertEqual(t, e, nil)
		util.AssertEqual(t, err.ErrCode, 0)
		util.AssertEqual(t, err.ErrMsg, "")
		util.AssertEqual(t, data.Bool, true)
		util.AssertEqual(t, data.String, "string")
		util.AssertEqual(t, data.Int, 1)
		util.AssertEqual(t, data.Float, 2.3)
	}

	{
		var (
			err  Err
			data Data
		)
		req, _ := http.NewRequest(http.MethodPost, ts.URL, nil)
		resp := do(req)
		e := resp.ConditionBindBody(condition, &err, &data).Error()

		util.AssertEqual(t, e, nil)
		util.AssertEqual(t, err.ErrCode, 1)
		util.AssertEqual(t, err.ErrMsg, "not ok")
		util.AssertEqual(t, data.Bool, false)
		util.AssertEqual(t, data.String, "")
		util.AssertEqual(t, data.Int, 0)
		util.AssertEqual(t, data.Float, 0.0)
	}
}

func do(req *http.Request) *Response {
	cli := &http.Client{}

	resp, err := cli.Do(req)

	return New(err, resp)
}
