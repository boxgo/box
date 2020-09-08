package wukong

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWithCTX(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	client := New(ts.URL)
	_, err := client.Get("/").WithCTX(ctx).End()
	if err != nil && !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatal(err)
	}
}

func TestSetHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, r.Header.Get("header_key1"), "header_value1")
		assertEqual(t, r.Header.Get("header_key2"), "header_value2")

		w.WriteHeader(200)
	}))
	defer ts.Close()

	client := New(ts.URL)
	_, err := client.Get("/").
		SetHeader("header_key1", "header_value1").
		SetHeader("header_key2", "header_value2").
		End()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddCookie(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			t.Error(err)
		}

		assertEqual(t, len(r.Cookies()), 1)
		assertEqual(t, cookie.Name, "session_id")
		assertEqual(t, cookie.Value, "abc")

		w.WriteHeader(200)
	}))
	defer ts.Close()

	client := New(ts.URL)
	_, err := client.Get("/").
		AddCookies(&http.Cookie{
			Name:  "session_id",
			Value: "abc",
		}).
		End()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, r.URL.String(), "/users/uid_123/friends/fid_456/images/1")

		w.WriteHeader(200)
	}))
	defer ts.Close()

	client := New(ts.URL)
	_, err := client.Get("/users/:uid/friends/:fid/images/:imgId").Param(map[string]interface{}{
		"uid":   "uid_123",
		"fid":   "fid_456",
		"imgId": 1,
	}).End()
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, r.URL.String(), "/?bool=true&float32=2.3&float64=4.5&int=1&int_array=1&int_array=2&interface_array=0&interface_array=true&interface_array=1.1&interface_array=string&string=string&string_array=a&string_array=b&uint=0")

		w.WriteHeader(200)
	}))
	defer ts.Close()

	client := New(ts.URL)
	_, err := client.Get("/").Query(map[string]interface{}{
		"string":          "string",
		"uint":            uint(0),
		"int":             1,
		"float32":         float32(2.3),
		"float64":         float32(4.5),
		"bool":            true,
		"string_array":    []string{"a", "b"},
		"int_array":       []int{1, 2},
		"interface_array": []interface{}{0, true, 1.1, "string"},
	}).End()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendMap(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, string(data), `{"fid":"fid_456","imgId":1,"uid":"uid_123"}`)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	client := New(ts.URL)
	_, err := client.Post("/").Send(map[string]interface{}{
		"uid":   "uid_123",
		"fid":   "fid_456",
		"imgId": 1,
	}).End()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendStruct(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, string(data), `{"uid":"uid_123","fid":"fid_456","imgId":1}`)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	td := struct {
		UId   string `json:"uid"`
		FId   string `json:"fid"`
		ImgId int    `json:"imgId"`
	}{
		UId:   "uid_123",
		FId:   "fid_456",
		ImgId: 1,
	}

	client := New(ts.URL)
	_, err := client.Post("/").Send(td).End()
	if err != nil {
		t.Fatal(err)
	}
}
