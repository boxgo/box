package ginlog_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boxgo/box/v2/logger"
	ginserver2 "github.com/boxgo/box/v2/server/ginserver"
	"github.com/boxgo/box/v2/server/ginserver/mid/ginlog"
	"github.com/gin-gonic/gin"
)

func Example() {
	ginserver2.Use(ginlog.Logger())
	ginserver2.GET("/ping", func(ctx *gin.Context) {
		ctx.Data(200, "text/plain", []byte("pong"))
	})

	if err := ginserver2.Run(); err != nil {
		logger.Fatal(err)
	}
}

// BenchmarkGin-8             	26134183	        49.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGin(b *testing.B) {
	server := ginserver2.StdConfig("default").Build()
	server.Use(func(ctx *gin.Context) {
		ctx.Next()
	})
	server.GET("/ping", func(ctx *gin.Context) {})

	req, _ := http.NewRequest("GET", "/ping", nil)
	recorder := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.ServeHTTP(recorder, req)
	}
}

// BenchmarkGinLog_Logger-8   	  467565	      2739 ns/op	    1097 B/op	      29 allocs/op
func BenchmarkGinLog_Logger(b *testing.B) {
	server := ginserver2.StdConfig("default").Build()
	server.Use(ginlog.Logger())
	server.GET("/ping", func(ctx *gin.Context) {})

	req, _ := http.NewRequest("GET", "/ping", nil)
	recorder := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.ServeHTTP(recorder, req)
	}
}
