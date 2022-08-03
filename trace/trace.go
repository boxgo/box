package trace

type (
	// Trace is application's information
	Trace struct {
		config *Config
	}
)

var (
	Default = StdConfig().Build()
)

func newTrace(cfg *Config) *Trace {
	return &Trace{
		config: cfg,
	}
}

// ID return application trace uid
func ID() string {
	return Default.config.TraceUID
}

// ReqID return application trace request id
func ReqID() string {
	return Default.config.TraceReqID
}

// SpanID return application trace span id
func SpanID() string {
	return Default.config.TraceSpanID
}

// BizID return application trace biz id
func BizID() string {
	return Default.config.TraceBizID
}
