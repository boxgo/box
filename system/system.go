package system

import (
	"os"
	"time"

	"github.com/boxgo/box/v2/util/netutil"
)

var (
	startAt  time.Time // Application start time.
	hostname string    // Runtime host's hostname.
	ip       string    // Runtime host's ip.
)

func init() {
	startAt = time.Now()
	hostname, _ = os.Hostname()
	ip, _ = netutil.GetGlobalUnicastIP()

	if ip == "" {
		ip = "127.0.0.1"
	}
	if hostname == "" {
		hostname = "localhost"
	}
}

// IP return host ip
func IP() string {
	return ip
}

// Hostname return hostname
func Hostname() string {
	return hostname
}

// StartAt return hostname
func StartAt() time.Time {
	return startAt
}
