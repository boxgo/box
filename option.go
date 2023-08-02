package box

import "github.com/boxgo/box/v2/log"

type (
	// Options new box options
	Options struct {
		id              string
		name            string
		version         string
		namespace       string
		metadata        map[string]string
		startupTimeout  int
		shutdownTimeout int
		autoMaxProcs    *bool
		log             log.Logger
		boxes           []Box
	}

	// Option setter
	Option func(ops *Options)
)

// WithId with service id.
func WithId(id string) Option {
	return func(o *Options) { o.id = id }
}

// WithName with service name.
func WithName(name string) Option {
	return func(o *Options) { o.name = name }
}

// WithVersion with service version.
func WithVersion(version string) Option {
	return func(o *Options) { o.version = version }
}

// WithNameSpace with service namespace.
func WithNameSpace(namespace string) Option {
	return func(o *Options) { o.namespace = namespace }
}

// WithMeta with service metadata.
func WithMeta(meta map[string]string) Option {
	return func(ops *Options) { ops.metadata = meta }
}

// WithLog with service id.
func WithLog(l log.Logger) Option {
	return func(o *Options) { o.log = l }
}

// WithStartupTimeout app startup timeout
func WithStartupTimeout(timeout int) Option {
	return func(ops *Options) { ops.startupTimeout = timeout }
}

// WithShutdownTimeout app shutdown timeout
func WithShutdownTimeout(timeout int) Option {
	return func(ops *Options) { ops.shutdownTimeout = timeout }
}

// WithAutoMaxProcs auto set `P`
func WithAutoMaxProcs(autoMaxProcs bool) Option {
	return func(ops *Options) { ops.autoMaxProcs = &autoMaxProcs }
}

// WithBoxes set boxes
func WithBoxes(boxes ...Box) Option {
	return func(ops *Options) { ops.boxes = boxes }
}
