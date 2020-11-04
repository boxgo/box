package field

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type (
	Fields []Field

	Config interface {
		Path() string
	}
)

func (fs Fields) Len() int {
	return len(fs)
}

func (fs Fields) Less(i, j int) bool {
	return fs[i].String() < fs[j].String()
}

func (fs Fields) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs *Fields) Sort() *Fields {
	sort.Sort(fs)

	return fs
}

func (fs *Fields) Table() string {
	fs.Sort()

	builder := &strings.Builder{}
	table := tablewriter.NewWriter(builder)
	table.SetHeader([]string{"Name", "Path", "Type", "Default", "Description"})
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetAutoFormatHeaders(false)
	table.SetRowLine(true)

	data := make([][]string, len(*fs))
	for idx, f := range *fs {
		data[idx] = f.Row()
	}

	table.AppendBulk(data)
	table.Render()

	return builder.String()
}

func (fs *Fields) Parse(c Config) *Fields {
	fs.parse(c.Path(), "", c)

	return fs
}

func (fs *Fields) parse(path, sub string, c interface{}) {
	var (
		n   = 0
		typ reflect.Type
		val reflect.Value
	)

	if reflect.TypeOf(c).Kind() == reflect.Ptr {
		typ = reflect.TypeOf(c).Elem()
		val = reflect.ValueOf(c).Elem()
		n = typ.NumField()
	} else {
		typ = reflect.TypeOf(c)
		val = reflect.ValueOf(c)
		n = typ.NumField()
	}

	for i := 0; i < n; i++ {
		if !val.Field(i).CanInterface() {
			continue
		}

		var (
			typField = typ.Field(i)
			valField = val.Field(i)
			tagField = typField.Tag
			name     = tagField.Get("config")
			desc     = tagField.Get("desc")
			def      = valField.Interface()
		)

		fn := func() {
			if sub != "" {
				name = sub + "." + name
			}

			*fs = append(*fs, Field{
				Path:    path,
				Name:    name,
				Type:    fmt.Sprintf("%T", def),
				Desc:    desc,
				Default: def,
			})
		}

		if typField.Type.Kind() == reflect.Struct {
			fs.parse(path, name, def)

			if !hasExport(valField) {
				fn()
			}
		} else if typField.Type.Kind() == reflect.Ptr && typField.Type.Elem().Kind() == reflect.Struct {
			fs.parse(path, name, def)

			if !hasExport(valField.Elem()) {
				fn()
			}
		} else {
			fn()
		}
	}
}

func hasExport(val reflect.Value) bool {
	n := val.NumField()
	for i := 0; i < n; i++ {
		if val.Field(i).CanInterface() {
			return true
		}
	}

	return false
}
