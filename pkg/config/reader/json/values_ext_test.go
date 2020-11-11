package json

import (
	"reflect"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
)

type (
	testcase struct {
		path []string
		typ  string
		def  interface{}
		exp  interface{}
	}

	fooStruct struct {
		String          string            `config:"string"`
		StringSlice     []string          `config:"string_slice"`
		StringMapString map[string]string `config:"string_map_string"`
		Int             int               `config:"int"`
		Uint            uint              `config:"uint"`
		Bool            bool              `config:"bool"`
		Float64         float64           `config:"float64"`
		Duration        time.Duration     `config:"duration"`
		Duration1       time.Duration     `config:"duration1"`
	}
)

var (
	fooBytes = []byte(`
		{
			"string": "string",
			"string_slice": ["1", "2"],
			"string_map_string": {
				"key1": "1",
				"key2": "2"
			},
			"int": -999,
			"uint": 999,
			"bool": true,
			"float64": 1.999999,
			"duration": "1h",
			"duration1": 1
		}
	`)
)

func test(t *testing.T, values reader.Values, ts ...testcase) {
	for _, it := range ts {
		v := values.Get(it.path...)
		var data interface{}

		pass := false
		switch it.typ {
		case "bool":
			data = v.Bool(it.def.(bool))
			pass = data == it.exp
		case "int":
			data = v.Int(it.def.(int))
			pass = data == it.exp
		case "uint":
			data = v.Uint(it.def.(uint))
			pass = data == it.exp
		case "string":
			data = v.String(it.def.(string))
			pass = data == it.exp
		case "float64":
			data = v.Float64(it.def.(float64))
			pass = data == it.exp
		case "duration":
			data = v.Duration(it.def.(time.Duration))
			pass = data == it.exp
		case "stringslice":
			data = v.StringSlice(it.def.([]string))
			pass = reflect.DeepEqual(data, it.exp)
		case "stringmapstring":
			data = v.StringMap(it.def.(map[string]string))
			pass = reflect.DeepEqual(data, it.exp)
		case "bytes":
			data = v.Bytes()
			pass = reflect.DeepEqual(data, it.exp)
		}

		t.Logf("%#v got %#v for path %v", it.exp, data, it.path)
		if !pass {
			t.Fatalf("Expected %#v got %#v for path %v", it.exp, data, it.path)
		}
	}
}

func TestExtValuesGet(t *testing.T) {
	values, err := newValues(&source.ChangeSet{
		Data: fooBytes,
	})
	if err != nil {
		t.Fatal(err)
	}

	test(
		t,
		values,
		[]testcase{
			testcase{path: []string{"string"}, typ: "string", exp: "string", def: ""},
			testcase{path: []string{"string_slice"}, typ: "stringslice", exp: []string{"1", "2"}, def: []string{}},
			testcase{path: []string{"string_slice.0"}, typ: "string", exp: "1", def: ""},
			testcase{path: []string{"string_slice.1"}, typ: "string", exp: "2", def: ""},
			testcase{path: []string{"string_slice", "0"}, typ: "string", exp: "1", def: ""},
			testcase{path: []string{"string_slice", "1"}, typ: "string", exp: "2", def: ""},
			testcase{path: []string{"string_map_string"}, typ: "stringmapstring", exp: map[string]string{"key1": "1", "key2": "2"}, def: map[string]string{}},
			testcase{path: []string{"string_map_string.key1"}, typ: "string", exp: "1", def: ""},
			testcase{path: []string{"string_map_string.key2"}, typ: "string", exp: "2", def: ""},
			testcase{path: []string{"string_map_string", "key1"}, typ: "string", exp: "1", def: ""},
			testcase{path: []string{"string_map_string", "key2"}, typ: "string", exp: "2", def: ""},
			testcase{path: []string{"int"}, typ: "int", exp: -999, def: 0},
			testcase{path: []string{"uint"}, typ: "uint", exp: uint(999), def: uint(0)},
			testcase{path: []string{"bool"}, typ: "bool", exp: true, def: false},
			testcase{path: []string{"float64"}, typ: "float64", exp: 1.999999, def: 0.0},
			testcase{path: []string{"duration"}, typ: "duration", exp: time.Hour, def: time.Duration(0)},
			testcase{path: []string{"duration1"}, typ: "duration", exp: time.Second, def: time.Duration(0)},
		}...,
	)
}

func TestExtValuesScan(t *testing.T) {
	values, err := newValues(&source.ChangeSet{
		Data: fooBytes,
	})
	if err != nil {
		t.Fatal(err)
	}

	foo := &fooStruct{}
	if err := values.Scan(foo); err != nil {
		t.Fatal(err)
	}

	test(
		t,
		values,
		[]testcase{
			testcase{path: []string{"string"}, typ: "string", exp: foo.String, def: ""},
			testcase{path: []string{"string_slice"}, typ: "stringslice", exp: foo.StringSlice, def: []string{}},
			testcase{path: []string{"string_slice.0"}, typ: "string", exp: foo.StringSlice[0], def: ""},
			testcase{path: []string{"string_slice.1"}, typ: "string", exp: foo.StringSlice[1], def: ""},
			testcase{path: []string{"string_slice", "0"}, typ: "string", exp: foo.StringSlice[0], def: ""},
			testcase{path: []string{"string_slice", "1"}, typ: "string", exp: foo.StringSlice[1], def: ""},
			testcase{path: []string{"string_map_string"}, typ: "stringmapstring", exp: foo.StringMapString, def: map[string]string{}},
			testcase{path: []string{"string_map_string.key1"}, typ: "string", exp: foo.StringMapString["key1"], def: ""},
			testcase{path: []string{"string_map_string.key2"}, typ: "string", exp: foo.StringMapString["key2"], def: ""},
			testcase{path: []string{"string_map_string", "key1"}, typ: "string", exp: foo.StringMapString["key1"], def: ""},
			testcase{path: []string{"string_map_string", "key2"}, typ: "string", exp: foo.StringMapString["key2"], def: ""},
			testcase{path: []string{"int"}, typ: "int", exp: foo.Int, def: 0},
			testcase{path: []string{"uint"}, typ: "uint", exp: foo.Uint, def: uint(0)},
			testcase{path: []string{"bool"}, typ: "bool", exp: foo.Bool, def: false},
			testcase{path: []string{"float64"}, typ: "float64", exp: foo.Float64, def: 0.0},
			testcase{path: []string{"duration"}, typ: "duration", exp: foo.Duration, def: time.Duration(0)},
			testcase{path: []string{"duration1"}, typ: "duration", exp: foo.Duration1, def: time.Duration(0)},
		}...,
	)
}

func TestExtValuesMap(t *testing.T) {
	values, err := newValues(&source.ChangeSet{
		Data: fooBytes,
	})
	if err != nil {
		t.Fatal(err)
	}

	foo := values.Map()

	test(
		t,
		values,
		[]testcase{
			testcase{path: []string{"string"}, typ: "string", exp: foo["string"], def: ""},
		}...,
	)
}

func TestExtValueScan(t *testing.T) {
	values, err := newValues(&source.ChangeSet{
		Data: fooBytes,
	})
	if err != nil {
		t.Fatal(err)
	}

	foo := &fooStruct{}
	values.Get("string").Scan(&foo.String)
	values.Get("string_slice").Scan(&foo.StringSlice)
	values.Get("string_map_string").Scan(&foo.StringMapString)
	values.Get("int").Scan(&foo.Int)
	values.Get("uint").Scan(&foo.Uint)
	values.Get("bool").Scan(&foo.Bool)
	values.Get("float64").Scan(&foo.Float64)
	values.Get("duration").Scan(&foo.Duration)
	values.Get("duration1").Scan(&foo.Duration1)

	test(
		t,
		values,
		[]testcase{
			testcase{path: []string{"string"}, typ: "string", exp: foo.String, def: ""},
			testcase{path: []string{"string_slice"}, typ: "stringslice", exp: foo.StringSlice, def: []string{}},
			testcase{path: []string{"string_slice.0"}, typ: "string", exp: foo.StringSlice[0], def: ""},
			testcase{path: []string{"string_slice.1"}, typ: "string", exp: foo.StringSlice[1], def: ""},
			testcase{path: []string{"string_slice", "0"}, typ: "string", exp: foo.StringSlice[0], def: ""},
			testcase{path: []string{"string_slice", "1"}, typ: "string", exp: foo.StringSlice[1], def: ""},
			testcase{path: []string{"string_map_string"}, typ: "stringmapstring", exp: foo.StringMapString, def: map[string]string{}},
			testcase{path: []string{"string_map_string.key1"}, typ: "string", exp: foo.StringMapString["key1"], def: ""},
			testcase{path: []string{"string_map_string.key2"}, typ: "string", exp: foo.StringMapString["key2"], def: ""},
			testcase{path: []string{"string_map_string", "key1"}, typ: "string", exp: foo.StringMapString["key1"], def: ""},
			testcase{path: []string{"string_map_string", "key2"}, typ: "string", exp: foo.StringMapString["key2"], def: ""},
			testcase{path: []string{"int"}, typ: "int", exp: foo.Int, def: 0},
			testcase{path: []string{"uint"}, typ: "uint", exp: foo.Uint, def: uint(0)},
			testcase{path: []string{"bool"}, typ: "bool", exp: foo.Bool, def: false},
			testcase{path: []string{"float64"}, typ: "float64", exp: foo.Float64, def: 0.0},
			testcase{path: []string{"duration"}, typ: "duration", exp: foo.Duration, def: time.Duration(0)},
			testcase{path: []string{"duration1"}, typ: "duration", exp: foo.Duration1, def: time.Duration(0)},
		}...,
	)
}

func (foo *fooStruct) Path() string {
	return ""
}
