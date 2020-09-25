package app

import (
	"fmt"
	"os"
	"time"

	"github.com/boxgo/box/pkg/util/netutil"
	"github.com/boxgo/box/pkg/util/strutil"
)

const (
	EnvName = "BOX_APP_NAME"
	EnvVer  = "BOX_APP_VERSION"
)

var (
	Name     string    // Application name. You can preset name when go build. And also, you can preset by env "BOX_APP_NAME".
	Version  string    // Application version. You can preset name when go build. And also, you can preset by env "BOX_APP_VERSION".
	StartAt  time.Time // Application start time.
	Hostname string    // Runtime host's hostname.
	IP       string    // Runtime host's ip.
)

func init() {
	Hostname, _ = os.Hostname()
	IP, _ = netutil.GetGlobalUnicastIP()
	StartAt = time.Now()

	if Name == "" {
		if env := os.Getenv(EnvName); env != "" {
			Name = env
		} else {
			Name = fmt.Sprintf("box_%s", strutil.RandString(6))
		}
	}

	if Version == "" {
		if env := os.Getenv(EnvVer); env != "" {
			Version = env
		} else {
			Version = "unknown"
		}
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
