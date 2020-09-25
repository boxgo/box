package field

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/boxgo/box/pkg/config/reader"
)

type (
	Field struct {
		Name      string
		Path      string
		Desc      string
		Immutable bool
		Def       interface{}
		val       reader.Value
		sync.Mutex
	}

	Fields []*Field
)

func New(immutable bool, name, path, desc string, def interface{}) *Field {
	return &Field{
		Immutable: immutable,
		Name:      name,
		Path:      path,
		Desc:      desc,
		Def:       def,
	}
}

func (f *Field) Row() []string {
	return []string{f.Name, f.Path, fmt.Sprintf("%t", f.Immutable), fmt.Sprintf("%v", f.Def), f.Desc}
}

func (f *Field) String() string {
	if f.Name != "" {
		return f.Name + "." + f.Path
	}

	return f.Path
}

func (f *Field) Paths() []string {
	return strings.Split(f.String(), ".")
}

func (f *Field) SetVal(val reader.Value) {
	f.Lock()
	defer f.Unlock()

	f.val = val
}

func (f *Field) Val() reader.Value {
	return f.val
}

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

func (fs *Fields) Append(fields ...*Field) *Fields {
	*fs = append(*fs, fields...)

	return fs
}
