// Package ginlog is gin server logger middleware.
package ginlog

import (
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	GinLog struct {
		cfg *Config
	}
)

func newGinLog(c *Config) *GinLog {
	return &GinLog{
		cfg: c,
	}
}

func (log *GinLog) Logger() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			fields     []interface{}
			start      = time.Now()
			method     = ctx.Request.Method
			path       = ctx.Request.URL.Path
			bodyWriter = newBodyWriter(ctx)
		)

		ctx.Writer = bodyWriter

		if log.cfg.RequestIP {
			fields = append(fields, "ip", ctx.ClientIP())
		}
		if log.cfg.RequestQuery {
			fields = append(fields, "query", ctx.Request.URL.RawQuery)
		}
		if log.cfg.RequestBody {
			fields = append(fields, "body", readBody(ctx))
		}
		if log.cfg.UserAgent {
			fields = append(fields, "user-agent", ctx.Request.UserAgent())
		}
		if log.cfg.RequestHeader {
			fields = append(fields, "header", ctx.Request.Header)
		}

		logger.Trace(ctx).Infow(fmt.Sprintf("http_server_req|%s|%s", method, path), fields...)

		ctx.Next()

		fields = append(fields, "status", ctx.Writer.Status(), "latency", time.Since(start))
		if log.cfg.ResponseBody {
			fields = append(fields, "resp", bodyWriter.body.String())
		}

		logger.Trace(ctx).Infow(fmt.Sprintf("http_server_resp|%s|%s", method, path), fields...)
	}
}
