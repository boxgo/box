package config

import (
	"sort"
)

type (
	Field struct {
		name string
		path string
		desc string
		def  interface{}
	}

	Fields []*Field
)

var (
	fieldBoxName     = NewField("box", "name", "box application name", "box")
	fieldTraceUid    = NewField("box", "trace.uid", "trace uid in context", "box.trace.uid")
	fieldTraceReqId  = NewField("box", "trace.reqId", "trace requestId in context", "box.trace.reqId")
	fieldTraceSpanId = NewField("box", "trace.spanId", "trace spanId in context", "box.trace.spanId")
	fieldTraceBizId  = NewField("box", "trace.bizId", "trace bizId in context", "box.trace.bizId")
)

func NewField(name, path, desc string, def interface{}) *Field {
	return &Field{
		name: name,
		path: path,
		desc: desc,
		def:  def,
	}
}

func (f Field) String() string {
	if f.name != "" {
		return f.name + "." + f.path
	}

	return f.path
}

func (fs Fields) Len() int {
	return len(fs)
}

func (fs Fields) Less(i, j int) bool {
	return fs[i].String() < fs[j].String()
}

func (fs Fields) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs *Fields) Sort() *Fields {
	sort.Sort(fs)

	return fs
}

func (fs *Fields) Append(fields ...*Field) *Fields {
	*fs = append(*fs, fields...)

	return fs
}
