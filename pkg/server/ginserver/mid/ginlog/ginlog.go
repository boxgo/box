// Package ginlog is gin server logger middleware.
package ginlog

import (
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/trace"
	"github.com/boxgo/box/pkg/util/strutil"
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
		if strutil.Contained(log.cfg.Skips, ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		var (
			fields []interface{}
			reqId  = ctx.GetHeader(trace.ReqID())
			start  = time.Now()
			method = ctx.Request.Method
			path   = ctx.Request.URL.Path
			writer *bodyWriter
		)

		if reqId == "" {
			reqId = strutil.RandomAlphanumeric(10)
		}
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
		if log.cfg.ResponseBody {
			writer = newBodyWriter(ctx)
			ctx.Writer = writer
		}

		ctx.Set(trace.ReqID(), reqId)
		ctx.Set(trace.BizID(), fmt.Sprintf("%s|%s", method, path))

		logger.Trace(ctx).Infow("http_server_req", fields...)

		ctx.Next()

		fields = append(fields, "status", ctx.Writer.Status(), "latency", time.Since(start))
		if log.cfg.ResponseBody && writer != nil {
			fields = append(fields, "resp", writer.body.String())
		}

		logger.Trace(ctx).Infow("http_server_rsp", fields...)
	}
}
