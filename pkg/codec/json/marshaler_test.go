package json_test

import (
	"testing"

	"github.com/boxgo/box/pkg/codec/json"
	"github.com/boxgo/box/pkg/codec/json/testdata/proto"
)

type (
	TestStruct struct {
		A int
		B string
		C bool
		D *string
	}
)

func TestJSONMarshal(t *testing.T) {
	testStruct := TestStruct{
		A: 1,
		B: "hello, world",
		C: true,
		D: nil,
	}

	codec := json.NewMarshaler()
	if data, err := codec.Marshal(testStruct); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%s", data)
	}
}

func TestProtoMarshal(t *testing.T) {
	testStruct := &proto.Ping{
		Msg:   "hi",
		Count: 100,
	}

	codec := json.NewMarshaler()
	if data, err := codec.Marshal(testStruct); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%s", data)
	}
}
