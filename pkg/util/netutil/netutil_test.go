package netutil

import (
	"testing"
)

func TestGetGlobalUnicastIP(t *testing.T) {
	ip, err := GetGlobalUnicastIP()
	t.Log(err, ip)

	if err != nil {
		t.Fatal(err)
	}
}
