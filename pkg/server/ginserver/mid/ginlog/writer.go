package ginlog

import (
	"bytes"
	"io/ioutil"

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

func readBody(ctx *gin.Context) []byte {
	if ctx.Request.Body == nil {
		return nil
	}

	body, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}
