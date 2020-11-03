package app

import (
	"fmt"
	"os"
	"time"

	"github.com/boxgo/box/pkg/util/netutil"
)

var (
	Name     string    // Application name. You can preset name when go build. Or it will generate a random name.
	Version  string    // Application version. You can preset name when go build. Or it will be setted to `unknown`.
	StartAt  time.Time // Application start time.
	Hostname string    // Runtime host's hostname.
	IP       string    // Runtime host's ip.
)

func init() {
	Hostname, _ = os.Hostname()
	IP, _ = netutil.GetGlobalUnicastIP()
	StartAt = time.Now()

	if Name == "" {
		Name = "box"
	}

	if Version == "" {
		Version = "unknown"
	}

	if IP == "" {
		IP = "127.0.0.1"
	}
	if Hostname == "" {
		Hostname = "localhost"
	}

}

// Summary get application runtime summary
func Summary() string {
	return fmt.Sprintf("Name: %s\nVersion: %s\nHostname: %s\nIP: %s\nStartAt: %s",
		Name, Version, Hostname, IP, StartAt.Format("2006-01-02 15:04:05"))
}
