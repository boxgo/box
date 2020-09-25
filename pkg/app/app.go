package app

import (
	"fmt"
	"os"
	"time"

	"github.com/boxgo/box/pkg/util/netutil"
	"github.com/boxgo/box/pkg/util/strutil"
)

const (
	envName = "BOX_APP_NAME"
	envVer  = "BOX_APP_VERSION"
)

var (
	Name     string
	Version  string
	Hostname string
	IP       string
	StartAt  time.Time
)

func init() {
	Hostname, _ = os.Hostname()
	IP, _ = netutil.GetGlobalUnicastIP()
	StartAt = time.Now()

	if Name == "" {
		if env := os.Getenv(envName); env != "" {
			Name = env
		} else {
			Name = fmt.Sprintf("box_%s", strutil.RandString(6))
		}
	}

	if Version == "" {
		if env := os.Getenv(envVer); env != "" {
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
