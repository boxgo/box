package runenv_test

import (
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/boxgo/box/v2/environment"
	"github.com/boxgo/box/v2/runenv"
)

func TestHostIP(t *testing.T) {
	ip := runenv.HostIP()

	reg := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	if !reg.MatchString(ip) {
		t.Fatalf("HostIP %s invalid", ip)
	}
}

func TestHostIP_ENV(t *testing.T) {
	hostIp := "192.168.1.100"

	os.Setenv(environment.HostIP, hostIp)
	defer os.Unsetenv(environment.HostIP)

	if ip := runenv.HostIP(); ip != hostIp {
		t.Fatalf("HostIP\nexpect:%s\nactual:%s", hostIp, ip)
	}
}

func TestPodIP(t *testing.T) {
	ip := runenv.PodIP()

	reg := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	if !reg.MatchString(ip) {
		t.Fatalf("PodIP %s invalid", ip)
	}
}

func TestPodIP_ENV(t *testing.T) {
	hostIp := "192.168.50.100"

	os.Setenv(environment.PodIP, hostIp)
	defer os.Unsetenv(environment.PodIP)

	if ip := runenv.PodIP(); ip != hostIp {
		t.Fatalf("PodIP\nexpect:%s\nactual:%s", hostIp, ip)
	}
}

func TestHostname(t *testing.T) {
	if runenv.Hostname() == "" {
		t.Fatalf("Hostname is empty")
	}
}

func TestStartAt(t *testing.T) {
	if runenv.StartAt().Before(time.Now().Add(-time.Second)) && runenv.StartAt().After(time.Now()) {
		t.Fatalf("Hostname is empty")
	}
}
