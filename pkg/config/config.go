// Package config is an interface for dynamic configuration.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/config/source/env"
	"github.com/boxgo/box/pkg/config/source/etcd"
	"github.com/boxgo/box/pkg/config/source/file"
	"github.com/boxgo/box/pkg/util"
)

type (
	// Configurator is an interface abstraction for dynamic configuration
	Configurator interface {
		// Load config sources
		Load(source ...source.Source) error
		// Force a source changeset sync
		Sync() error
		// Stop the config loader/watcher
		Close() error
		// Bytes get merged config data
		Bytes() []byte
		// SprintFields registered fields
		SprintFields() string
		// SprintTemplate through encoder
		SprintTemplate(encoder string) string
		// Watch, Mount, Getter value through field
		WatchMountGetter
	}

	// Configurable instance should implements this interface
	Config interface {
		Name() string
		Configure(mg WatchMountGetter)
	}

	WatchMountGetter interface {
		// Watch field change
		Watch(field *Field) (Watcher, error)
		// Mount fields to configurator
		Mount(fields ...*Field)
		// Get value through field
		Get(field *Field) reader.Value
		// GetString through field
		GetBool(field *Field) bool
		// GetInt through field
		GetInt(field *Field) int
		// GetUint through field
		GetUint(field *Field) uint
		// GetString through field
		GetString(field *Field) string
		// GetFloat64 through field
		GetFloat64(field *Field) float64
		// GetDuration through field
		GetDuration(field *Field) time.Duration
		// GetStringSlice through field
		GetStringSlice(field *Field) []string
		// GetStringMap through field
		GetStringMap(field *Field) map[string]string
	}
)

var (
	// Default Config Manager
	Default = NewConfig()
)

// NewSimple create a configurator with `file` and `env` source support.
// the priority is: `env` > `file`.
// `env` will use filename as prefix automatically.
func NewSimple(filePath string) Configurator {
	name := util.Filename(filePath)
	nameUpper := strings.ToUpper(name)

	return newConfig(
		WithSource(
			file.NewSource(
				file.WithPath(filePath),
			),
			env.NewSource(
				env.WithPrefix(nameUpper),
				env.WithStrippedPrefix(nameUpper),
			),
		),
	)
}

// NewClassic create a configurator with `file`, `env` and `etcd` source support.
// the priority is: `etcd` > `env` > `file`.
// `env` will use filename as prefix automatically.
// `etcd` key format: `/{fileName}/config`
func NewClassic(filePath, username, password, address string) Configurator {
	name := util.Filename(filePath)
	nameUpper := strings.ToUpper(name)

	return newConfig(
		WithSource(
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
		),
	)
}

// NewEtcd create a configurator with `etcd` source support.
// `etcd` key format: `/{prefix}/config`
func NewEtcd(prefix, username, password, address string) Configurator {
	return newConfig(
		WithSource(
			etcd.NewSource(
				etcd.WithAddress(address),
				etcd.Auth(username, password),
				etcd.WithPrefix(fmt.Sprintf("/%s/config", prefix)),
				etcd.StripPrefix(true),
			),
		),
	)
}

// NewConfig returns new config
func NewConfig(opts ...Option) Configurator {
	return newConfig(opts...)
}

// LoadFile is short hand for creating a file source and loading it
func LoadFile(path string) error {
	return Load(file.NewSource(
		file.WithPath(path),
	))
}

// Load config sources
func Load(source ...source.Source) error {
	return Default.Load(source...)
}

// Force a source changeset sync
func Sync() error {
	return Default.Sync()
}

// Watch a value for changes
func Watch(field *Field) (Watcher, error) {
	return Default.Watch(field)
}

// Stop the config loader/watcher
func Close() error {
	return Default.Close()
}

// Mount fields
func Mount(fields ...*Field) {
	Default.Mount(fields...)
}

// Get a value from the config
func Get(field *Field) reader.Value {
	return Default.Get(field)
}

// GetString through field
func GetBool(field *Field) bool {
	return Default.GetBool(field)
}

// GetInt through field
func GetInt(field *Field) int {
	return Default.GetInt(field)
}

// GetUint through field
func GetUint(field *Field) uint {
	return Default.GetUint(field)
}

// GetString through field
func GetString(field *Field) string {
	return Default.GetString(field)
}

// GetFloat64 through field
func GetFloat64(field *Field) float64 {
	return Default.GetFloat64(field)
}

// GetDuration through field
func GetDuration(field *Field) time.Duration {
	return Default.GetDuration(field)
}

// GetStringSlice through field
func GetStringSlice(field *Field) []string {
	return Default.GetStringSlice(field)
}

// GetStringMap through field
func GetStringMap(field *Field) map[string]string {
	return Default.GetStringMap(field)
}

// SprintFields registered fields
func SprintFields() string {
	return Default.SprintFields()
}

// SprintTemplate through encoder
func SprintTemplate(encoder string) string {
	return Default.SprintTemplate(encoder)
}
