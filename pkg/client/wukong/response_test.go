package wukong

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/util/testutil"
)

func TestResponseIsTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	resp := New(ts.URL).Get("/").WithCTX(ctx).End()

	testutil.ExpectEqual(t, true, resp.IsTimeout())
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

	resp := New(ts.URL).Get("/").WithCTX(ctx).End()

	testutil.ExpectEqual(t, true, resp.IsCancel())
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

	resp := New(ts.URL).Get("/").WithCTX(ctx).End().BindBody(d)

	testutil.ExpectEqual(t, resp.Error(), nil)
	testutil.ExpectEqual(t, d.String, "string")
	testutil.ExpectEqual(t, d.Int, 1)
	testutil.ExpectEqual(t, d.Float, 2.3)
	testutil.ExpectEqual(t, d.Bool, true)
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
		e := New(ts.URL).Get("/").End().ConditionBindBody(condition, &err, &data).Error()

		testutil.ExpectEqual(t, e, nil)
		testutil.ExpectEqual(t, err.ErrCode, 0)
		testutil.ExpectEqual(t, err.ErrMsg, "")
		testutil.ExpectEqual(t, data.Bool, true)
		testutil.ExpectEqual(t, data.String, "string")
		testutil.ExpectEqual(t, data.Int, 1)
		testutil.ExpectEqual(t, data.Float, 2.3)
	}

	{
		var (
			err  Err
			data Data
		)
		e := New(ts.URL).Post("/").End().ConditionBindBody(condition, &err, &data).Error()

		testutil.ExpectEqual(t, e, nil)
		testutil.ExpectEqual(t, err.ErrCode, 1)
		testutil.ExpectEqual(t, err.ErrMsg, "not ok")
		testutil.ExpectEqual(t, data.Bool, false)
		testutil.ExpectEqual(t, data.String, "")
		testutil.ExpectEqual(t, data.Int, 0)
		testutil.ExpectEqual(t, data.Float, 0.0)
	}
}
