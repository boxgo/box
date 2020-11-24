package testutil

import (
	"reflect"
	"testing"
)

func ExpectEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("\nexpect a == b\nactual a = %#v , b = %#v", a, b)
	} else {
		t.Logf("\nexpect a == b\nactual a = %#v , b = %#v", a, b)
	}
}
