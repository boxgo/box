package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/boxgo/box/pkg/config/loader"
	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/config/source/env"
	"github.com/boxgo/box/pkg/config/source/etcd"
	"github.com/boxgo/box/pkg/config/source/file"
	"github.com/boxgo/box/pkg/util"
)

type (
	Options struct {
		Loader loader.Loader
		Reader reader.Reader
		Source []source.Source

		// for alternative data
		Context context.Context
	}

	Option func(o *Options)
)

// WithLoader sets the loader for manager config
func WithLoader(l loader.Loader) Option {
	return func(o *Options) {
		o.Loader = l
	}
}

// WithSource appends a source to list of sources
func WithSource(s ...source.Source) Option {
	return func(o *Options) {
		o.Source = append(o.Source, s...)
	}
}

// WithReader sets the config reader
func WithReader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}

// WithSimpleSource create a option with `file` and `env` source support.
// the priority is: `env` > `file`.
// `env` will use filename as prefix automatically.
func WithSimpleSource(filePath string) Option {
	return WithSource(SimpleSource(filePath)...)
}

// WithClassicSource create a option with `file`, `env` and `etcd` source support.
// the priority is: `etcd` > `env` > `file`.
// `env` will use filename as prefix automatically.
// `etcd` key format: `/{fileName}/config`
func WithClassicSource(filePath, username, password, address string) Option {
	return WithSource(ClassicSource(filePath, username, password, address)...)
}

// WithEtcdSource create a option with `etcd` source support.
// `etcd` key format: `/{prefix}/config`
func WithEtcdSource(prefix, username, password, address string) Option {
	return WithSource(EtcdSource(prefix, username, password, address)...)
}

func SimpleSource(filePath string) []source.Source {
	name := util.Filename(filePath)
	nameUpper := strings.ToUpper(name)

	return []source.Source{
		file.NewSource(
			file.WithPath(filePath),
		),
		env.NewSource(
			env.WithPrefix(nameUpper),
			env.WithStrippedPrefix(nameUpper),
		),
	}
}

func ClassicSource(filePath, username, password, address string) []source.Source {
	name := util.Filename(filePath)
	nameUpper := strings.ToUpper(name)

	return []source.Source{
		file.NewSource(
			file.WithPath(filePath),
		),
		env.NewSource(
			env.WithPrefix(nameUpper),
			env.WithStrippedPrefix(nameUpper),
		),
		etcd.NewSource(
			etcd.WithAddress(address),
			etcd.Auth(username, password),
			etcd.WithPrefix(fmt.Sprintf("/%s/config", name)),
			etcd.StripPrefix(true),
		),
	}
}

func EtcdSource(prefix, username, password, address string) []source.Source {
	return []source.Source{
		etcd.NewSource(
			etcd.WithAddress(address),
			etcd.Auth(username, password),
			etcd.WithPrefix(fmt.Sprintf("/%s/config", prefix)),
			etcd.StripPrefix(true),
		),
	}
}
