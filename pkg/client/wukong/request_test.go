package wukong

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/util/testutil"
)

func TestWithTimeoutCTX(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	req := New(ts.URL).Get("/").WithCTX(ctx)
	resp := req.End()
	if err := resp.Error(); err != nil && !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatal(err)
	}

	t.Log(req.TraceInfo.ElapsedTime)
	t.Log(req.TraceInfo.ConnectElapsed)
	t.Log(req.TraceInfo.GetConnElapsed)
	t.Log(req.TraceInfo.DNSLookupElapsed)
}

func TestWithCancelCTX(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Millisecond * 500)
		cancel()
	}()

	resp := New(ts.URL).Get("/").WithCTX(ctx).End()
	if err := resp.Error(); err != nil && !strings.Contains(err.Error(), "context canceled") {
		t.Fatal(err)
	}
}

func TestSetHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.ExpectEqual(t, r.Header.Get("header_key1"), "header_value1")
		testutil.ExpectEqual(t, r.Header.Get("header_key2"), "header_value2")

		w.WriteHeader(200)
	}))
	defer ts.Close()

	resp := New(ts.URL).Get("/").
		SetHeader("header_key1", "header_value1").
		SetHeader("header_key2", "header_value2").
		End()
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestAddCookie(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			t.Error(err)
		}

		testutil.ExpectEqual(t, len(r.Cookies()), 1)
		testutil.ExpectEqual(t, cookie.Name, "session_id")
		testutil.ExpectEqual(t, cookie.Value, "abc")

		w.WriteHeader(200)
	}))
	defer ts.Close()

	resp := New(ts.URL).Get("/").
		AddCookies(&http.Cookie{
			Name:  "session_id",
			Value: "abc",
		}).
		End()
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.ExpectEqual(t, r.URL.String(), "/users/uid_123/friends/fid_456/images/1")

		w.WriteHeader(200)
	}))
	defer ts.Close()

	resp := New(ts.URL).Get("/users/:uid/friends/:fid/images/:imgId").
		Param(map[string]interface{}{
			"uid":   "uid_123",
			"fid":   "fid_456",
			"imgId": 1,
		}).End()
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.ExpectEqual(t, r.URL.String(), "/?bool=true&float32=2.3&float64=4.5&int=1&int_array=1&int_array=2&interface_array=0&interface_array=true&interface_array=1.1&interface_array=string&string=string&string_array=a&string_array=b&uint=0")

		w.WriteHeader(200)
	}))
	defer ts.Close()

	resp := New(ts.URL).Get("/").Query(map[string]interface{}{
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
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestSendJsonMap(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		testutil.ExpectEqual(t, string(data), `{"fid":"fid_456","imgId":1,"uid":"uid_123"}`)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	resp := New(ts.URL).Post("/").Send(map[string]interface{}{
		"uid":   "uid_123",
		"fid":   "fid_456",
		"imgId": 1,
	}).End()
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestSendJsonStruct(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		testutil.ExpectEqual(t, string(data), `{"uid":"uid_123","fid":"fid_456","imgId":1}`)

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

	resp := New(ts.URL).Post("/").Send(td).End()
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestSendXmlStruct(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		testutil.ExpectEqual(t, string(data), `<Test><uid>uid_123</uid><fid>fid_456</fid><imgId>1</imgId></Test>`)

		w.WriteHeader(200)
	}))
	defer ts.Close()

	type Test struct {
		UId   string `xml:"uid"`
		FId   string `xml:"fid"`
		ImgId int    `xml:"imgId"`
	}

	td := Test{
		UId:   "uid_123",
		FId:   "fid_456",
		ImgId: 1,
	}

	resp := New(ts.URL).Post("/").Type(MimeTypeXML).Send(td).End()
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestSendFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		if fileHeader.Size != 13 {
			t.Fatal("文件大小错误")
		}
		if fileHeader.Filename != "upload.txt" {
			t.Fatal("文件名错误")
		}
		if data, _ := ioutil.ReadAll(file); string(data) != "hello world!\n" {
			t.Fatal("文件内容错误")
		}

		w.WriteHeader(200)
	}))
	defer ts.Close()

	pwd, _ := os.Getwd()
	fp := filepath.Join(pwd, "testdata/upload.txt")
	resp := New(ts.URL).Post("/").SendFile("file", fp).End()
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestSendFileReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, fileHeader, err := r.FormFile("file1")
		if err != nil {
			t.Fatal(err)
		}

		defer file.Close()

		if fileHeader.Size != 12 {
			t.Fatal("文件大小错误")
		}
		if fileHeader.Filename != "upload.txt" {
			t.Fatal("文件名错误")
		}
		if data, _ := ioutil.ReadAll(file); string(data) != "hello world!" {
			t.Fatal("文件内容错误")
		}

		w.WriteHeader(200)
	}))
	defer ts.Close()

	pwd, _ := os.Getwd()
	fp := filepath.Join(pwd, "testdata/upload.txt")
	resp := New(ts.URL).Post("/").SendFileReader("file1", fp, bytes.NewBuffer([]byte("hello world!"))).End()
	if err := resp.Error(); err != nil {
		t.Fatal(err)
	}
}
