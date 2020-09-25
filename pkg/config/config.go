// Package config is an interface for dynamic configuration.
package config

import (
	"time"

	"github.com/boxgo/box/pkg/config/field"
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
		// Watch, Mount, Getter value through field
		SubConfigurator
	}

	SubConfigurator interface {
		// Watch field change
		Watch(field *field.Field) (Watcher, error)
		// Mount fields to configurator
		Mount(fields ...*field.Field)
		// Get value through field
		Get(field *field.Field) reader.Value
		// GetString through field
		GetBool(field *field.Field) bool
		// GetInt through field
		GetInt(field *field.Field) int
		// GetUint through field
		GetUint(field *field.Field) uint
		// GetString through field
		GetString(field *field.Field) string
		// GetFloat64 through field
		GetFloat64(field *field.Field) float64
		// GetDuration through field
		GetDuration(field *field.Field) time.Duration
		// GetStringSlice through field
		GetStringSlice(field *field.Field) []string
		// GetStringMap through field
		GetStringMap(field *field.Field) map[string]string
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
func Watch(field *field.Field) (Watcher, error) {
	return Default.Watch(field)
}

// Stop the config loader/watcher
func Close() error {
	return Default.Close()
}

// Mount fields
func Mount(fields ...*field.Field) {
	Default.Mount(fields...)
}

// Get a value from the config
func Get(field *field.Field) reader.Value {
	return Default.Get(field)
}

// GetString through field
func GetBool(field *field.Field) bool {
	return Default.GetBool(field)
}

// GetInt through field
func GetInt(field *field.Field) int {
	return Default.GetInt(field)
}

// GetUint through field
func GetUint(field *field.Field) uint {
	return Default.GetUint(field)
}

// GetString through field
func GetString(field *field.Field) string {
	return Default.GetString(field)
}

// GetFloat64 through field
func GetFloat64(field *field.Field) float64 {
	return Default.GetFloat64(field)
}

// GetDuration through field
func GetDuration(field *field.Field) time.Duration {
	return Default.GetDuration(field)
}

// GetStringSlice through field
func GetStringSlice(field *field.Field) []string {
	return Default.GetStringSlice(field)
}

// GetStringMap through field
func GetStringMap(field *field.Field) map[string]string {
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
