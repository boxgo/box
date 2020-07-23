// Package config is an interface for dynamic configuration.
package config

import (
	"time"

	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/source"
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
		SubConfigurator
	}

	SubConfigurator interface {
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
		// GetBoxName path: box.name
		GetBoxName() string
		// GetTraceUid path: box.trace.uid
		GetTraceUid() string
		// GetTraceReqId path: box.trace.reqid
		GetTraceReqId() string
		// GetTraceBizId path: box.trace.bizid
		GetTraceBizId() string
		// GetTraceSpanId path: box.trace.spanid
		GetTraceSpanId() string
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
	return newConfig(WithSimpleSource(filePath))
}

// NewClassic create a configurator with `file`, `env` and `etcd` source support.
// the priority is: `etcd` > `env` > `file`.
// `env` will use filename as prefix automatically.
// `etcd` key format: `/{fileName}/config`
func NewClassic(filePath, username, password, address string) Configurator {
	return newConfig(WithClassicSource(filePath, username, password, address))
}

// NewEtcd create a configurator with `etcd` source support.
// `etcd` key format: `/{prefix}/config`
func NewEtcd(prefix, username, password, address string) Configurator {
	return newConfig(WithEtcdSource(prefix, username, password, address))
}

// NewConfig returns new config
func NewConfig(opts ...Option) Configurator {
	return newConfig(opts...)
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

// GetBoxName path: box.name
func GetBoxName() string {
	return Default.GetBoxName()
}

// GetTraceUid path: box.trace.uid
func GetTraceUid() string {
	return Default.GetTraceUid()
}

// GetTraceReqId path: box.trace.reqid
func GetTraceReqId() string {
	return Default.GetTraceReqId()
}

// GetTraceBizId path: box.trace.bizid
func GetTraceBizId() string {
	return Default.GetTraceBizId()
}

// GetTraceSpanId path: box.trace.spanid
func GetTraceSpanId() string {
	return Default.GetTraceSpanId()
}

// SprintFields registered fields
func SprintFields() string {
	return Default.SprintFields()
}

// SprintTemplate through encoder
func SprintTemplate(encoder string) string {
	return Default.SprintTemplate(encoder)
}
