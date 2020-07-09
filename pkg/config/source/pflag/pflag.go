package pflag

import (
	"errors"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/config/source"
	"github.com/imdario/mergo"
	"github.com/spf13/pflag"
)

type (
	pflagsrc struct {
		opts source.Options
	}
)

func (pfs *pflagsrc) String() string {
	return "pflag"
}

func (pfs *pflagsrc) Watch() (source.Watcher, error) {
	return source.NewNoopWatcher()
}

func (pfs *pflagsrc) Read() (*source.ChangeSet, error) {
	if !pflag.Parsed() {
		return nil, errors.New("flags not parsed")
	}

	var changes map[string]interface{}

	visitFn := func(f *pflag.Flag) {
		n := strings.ToLower(f.Name)
		keys := strings.FieldsFunc(n, split)
		reverse(keys)

		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				var v interface{}
				switch f.Value.Type() {
				case "bool":
					v = f.Value
				case "float32", "float64":
					v = f.Value
				case "int", "int8", "int16", "int32", "int64":
					v = f.Value
				case "uint", "uint8", "uint16", "uint32", "uint64":
					v = f.Value
				case "stringSlice":
					v, _ = pflag.CommandLine.GetStringSlice(f.Name)
				case "intSlice":
					v, _ = pflag.CommandLine.GetIntSlice(f.Name)
				case "int32Slice":
					v, _ = pflag.CommandLine.GetInt32Slice(f.Name)
				case "int64Slice":
					v, _ = pflag.CommandLine.GetInt64Slice(f.Name)
				case "uintSlice":
					v, _ = pflag.CommandLine.GetUintSlice(f.Name)
				case "boolSlice":
					v, _ = pflag.CommandLine.GetBoolSlice(f.Name)
				case "durationSlice":
					v, _ = pflag.CommandLine.GetDurationSlice(f.Name)
				default:
					v = f.Value
				}

				tmp[k] = v
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		mergo.Map(&changes, tmp) // need to sort error handling
		return
	}

	unset, ok := pfs.opts.Context.Value(includeUnsetKey{}).(bool)
	if ok && unset {
		pflag.VisitAll(visitFn)
	} else {
		pflag.Visit(visitFn)
	}

	b, err := pfs.opts.Encoder.Encode(changes)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    pfs.opts.Encoder.String(),
		Data:      b,
		Timestamp: time.Now(),
		Source:    pfs.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

// NewSource returns a config source for integrating parsed pflags.
// Hyphens are delimiters for nesting, and all keys are lowercased.
//
// Example:
//      dbhost := flag.String("database-host", "localhost", "the db host name")
//
//      {
//          "database": {
//              "host": "localhost"
//          }
//      }
func NewSource(opts ...source.Option) source.Source {
	return &pflagsrc{opts: source.NewOptions(opts...)}
}

func split(r rune) bool {
	return r == '-' || r == '_' || r == '.'
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}
