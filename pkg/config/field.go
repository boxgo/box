package config

type Field struct {
	name string
	path string
	desc string
	def  interface{}
}

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
