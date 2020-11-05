package system

import (
	"os"
	"time"

	"github.com/boxgo/box/pkg/util/netutil"
)

type (
	// System is application's information
	System struct {
		config        *Config
		runtimeStatus RuntimeStatus
	}

	// RuntimeStatus runtime status
	RuntimeStatus struct {
		StartAt  time.Time // Application start time.
		Hostname string    // Runtime host's hostname.
		IP       string    // Runtime host's ip.
	}
)

var (
	ip       string
	hostname string
	startAt  time.Time
	Default  *System
)

func init() {
	hostname, _ = os.Hostname()
	ip, _ = netutil.GetGlobalUnicastIP()
	startAt = time.Now()

	if ip == "" {
		ip = "127.0.0.1"
	}
	if hostname == "" {
		hostname = "localhost"
	}

	Default = StdConfig().Build()
}

func newSystem(cfg *Config) *System {
	return &System{
		config: cfg,
		runtimeStatus: RuntimeStatus{
			StartAt:  startAt,
			Hostname: hostname,
			IP:       ip,
		},
	}
}

// Runtime return application runtime status.
func Runtime() RuntimeStatus {
	return Default.runtimeStatus
}

// IP return host ip
func IP() string {
	return Runtime().IP
}

// Hostname return hostname
func Hostname() string {
	return Runtime().Hostname
}

// StartAt return hostname
func StartAt() time.Time {
	return Runtime().StartAt
}

// ServiceName return application name
func ServiceName() string {
	return Default.config.Name
}

// ServiceVersion return application version
func ServiceVersion() string {
	return Default.config.Version
}

// TraceUID return application trace uid
func TraceUID() string {
	return Default.config.TraceUID
}

// TraceReqID return application trace request id
func TraceReqID() string {
	return Default.config.TraceReqID
}

// TraceSpanID return application trace span id
func TraceSpanID() string {
	return Default.config.TraceSpanID
}

// TraceBizID return application trace biz id
func TraceBizID() string {
	return Default.config.TraceBizID
}
