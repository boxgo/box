package errcode

import (
	"testing"
)

func TestErrCode(t *testing.T) {
	if ErrrModDBReadTimeout.Code() != 100100001 {
		t.Fatalf("errcode should be %d", ErrrModDBReadTimeout.Code())
	}
}
