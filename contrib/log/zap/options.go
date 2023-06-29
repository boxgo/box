package zap

type (
	Options struct {
		traceFields []string
		maskRules   MaskRules
	}

	Option func(o *Options)
)

func WithTraceFields(fields []string) Option {
	return func(o *Options) {
		o.traceFields = fields
	}
}

func WithMaskRules(rules MaskRules) Option {
	return func(o *Options) {
		o.maskRules = rules
	}
}
