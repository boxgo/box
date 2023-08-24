package utils_test

import (
	"regexp"
	"testing"

	"github.com/boxgo/box/v2/internal/utils"
)

func TestGlobalIP4(t *testing.T) {
	ip, err := utils.GlobalIP4()
	if err != nil {
		t.Fatal(err)
	}

	reg := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	if !reg.MatchString(ip) {
		t.Fatal("GlobalIP4.Error", ip)
	}
}
