package wukong_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boxgo/box/pkg/client/wukong"
	"github.com/boxgo/box/pkg/util/testutil"
)

func send(b *testing.B, client *wukong.WuKong) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"string":"string", "int": 1, "float": 2.3, "bool": true}`))
	}))
	defer ts.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
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
			err := client.Get(ts.URL).WithCTX(context.TODO()).End().
				BindStatusCode(&statusCode).
				BindStatus(&statusMsg).
				BindHeader(&header).
				BindIsTimeout(&isTimeout).
				BindIsCancel(&isCancel).
				BindBody(&bodyData).
				Error()

			testutil.ExpectEqualB(b, err, nil)
			testutil.ExpectEqualB(b, statusCode, 200)
			testutil.ExpectEqualB(b, statusMsg, "200 OK")
			testutil.ExpectEqualB(b, isTimeout, false)
			testutil.ExpectEqualB(b, isCancel, false)
		}
	})
}

func BenchmarkSend1Client(b *testing.B) {
	send(b, wukong.New("").SetClientCount(1).Logger(wukong.LoggerDisable))
}

func BenchmarkSend5Client(b *testing.B) {
	send(b, wukong.New("").SetClientCount(5).Logger(wukong.LoggerDisable))
}

func BenchmarkSend8Client(b *testing.B) {
	send(b, wukong.New("").SetClientCount(8).Logger(wukong.LoggerDisable))
}

func BenchmarkSend10Client(b *testing.B) {
	send(b, wukong.New("").SetClientCount(10).Logger(wukong.LoggerDisable))
}

func BenchmarkSend20Client(b *testing.B) {
	send(b, wukong.New("").SetClientCount(20).Logger(wukong.LoggerDisable))
}

func BenchmarkSend40Client(b *testing.B) {
	send(b, wukong.New("").SetClientCount(40).Logger(wukong.LoggerDisable))
}

func BenchmarkSend80Client(b *testing.B) {
	send(b, wukong.New("").SetClientCount(80).Logger(wukong.LoggerDisable))
}

func BenchmarkSend100Client(b *testing.B) {
	send(b, wukong.New("").SetClientCount(100).Logger(wukong.LoggerDisable))
}
