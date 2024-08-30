package testutil

import (
	"reflect"
	"testing"
)

func ExpectEqual(t *testing.T, left, right interface{}) {
	if !reflect.DeepEqual(left, right) {
		t.Errorf("\nexpect a == b\nactual\na = %#v\nb = %#v", left, right)
	}
}

func ExpectEqualB(b *testing.B, left, right interface{}) {
	if !reflect.DeepEqual(left, right) {
		b.Errorf("\nexpect a == b\nactual\na = %#v\nb = %#v", left, right)
	}
}
