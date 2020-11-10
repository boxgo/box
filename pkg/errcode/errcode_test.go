package errcode

import (
	"testing"
)

func TestErrCode(t *testing.T) {
	t.Log(ErrrModDBDeleteTimeout1.Code(), ErrrModDBDeleteTimeout1.Message())
}
