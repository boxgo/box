package runenv

import (
	"os"
	"time"

	"github.com/boxgo/box/v2/environment"
	"github.com/boxgo/box/v2/internal/utils"
)

var (
	startAt time.Time
)

func init() {
	startAt = time.Now()
}

// IP retrieves ip
func IP() string {
	ip, _ := utils.GlobalIP4()

	if ip == "" {
		ip = "127.0.0.1"
	}

	return ip
}

// HostIP retrieves BOX_HOST_IP environment if exists, otherwise IP
func HostIP() string {
	if ip := os.Getenv(environment.HostIP); ip != "" {
		return ip
	}

	return IP()
}

// PodIP retrieves BOX_POD_IP environment if exists, otherwise IP
func PodIP() string {
	if ip := os.Getenv(environment.PodIP); ip != "" {
		return ip
	}

	return IP()
}

// Hostname retrieves hostname
func Hostname() string {
	hostname, _ := os.Hostname()

	if hostname == "" {
		hostname = "localhost"
	}

	return hostname
}

// StartAt retrieves application start time.
func StartAt() time.Time {
	return startAt
}
