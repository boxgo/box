package ginlog

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func newBodyWriter(ctx *gin.Context) *bodyWriter {
	return &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
}

func readBody(ctx *gin.Context) string {
	if ctx.Request.Body == nil {
		return ""
	}

	if !strings.Contains(ctx.ContentType(), "application/json") {
		return "<non-json ignored>"
	}

	body, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return string(body)
}
