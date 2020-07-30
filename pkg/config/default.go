package config

import (
	"sync"
	"time"

	"github.com/boxgo/box/pkg/config/loader"
	"github.com/boxgo/box/pkg/config/loader/memory"
	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/reader/json"
	"github.com/boxgo/box/pkg/config/source"
)

type (
	config struct {
		exit chan bool
		opts Options

		sync.RWMutex
		// the current snapshot
		snap *loader.Snapshot
		// the current values
		vals reader.Values
		// system registered fields
		sysFields Fields
		// user registered fields
		userFields Fields
	}
)

func newConfig(opts ...Option) Configurator {
	options := Options{
		Loader: memory.NewLoader(),
		Reader: json.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	options.Loader.Load(options.Source...)
	snap, _ := options.Loader.Snapshot()
	vals, _ := options.Reader.Values(snap.ChangeSet)

	c := &config{
		exit:       make(chan bool),
		opts:       options,
		snap:       snap,
		vals:       vals,
		sysFields:  make(Fields, 0),
		userFields: make(Fields, 0),
	}

	c.MountSystem(
		fieldBoxName,
		fieldTraceUid,
		fieldTraceReqId,
		fieldTraceSpanId,
		fieldTraceBizId,
	)

	go c.run()

	return c
}

func (c *config) run() {
	watch := func(w loader.Watcher) error {
		for {
			// get changeset
			snap, err := w.Next()
			if err != nil {
				return err
			}

			c.Lock()

			// save
			c.snap = snap

			// set values
			c.vals, _ = c.opts.Reader.Values(snap.ChangeSet)

			c.Unlock()
		}
	}

	for {
		w, err := c.opts.Loader.Watch()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		done := make(chan bool)

		// the stop watch func
		go func() {
			select {
			case <-done:
			case <-c.exit:
			}
			w.Stop()
		}()

		// block watch
		if err := watch(w); err != nil {
			// do something better
			time.Sleep(time.Second)
		}

		// close done chan
		close(done)

		// if the config is closed exit
		select {
		case <-c.exit:
			return
		default:
		}
	}
}

// Load config sources
func (c *config) Load(sources ...source.Source) error {
	if err := c.opts.Loader.Load(sources...); err != nil {
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	vals, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.vals = vals

	return nil
}

// sync loads all the sources, calls the parser and updates the config
func (c *config) Sync() error {
	if err := c.opts.Loader.Sync(); err != nil {
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	vals, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.vals = vals

	return nil
}

// Watch a value for changes
func (c *config) Watch(field *Field) (Watcher, error) {
	value := c.Get(field)
	path := field2path(field)

	w, err := c.opts.Loader.Watch(path...)
	if err != nil {
		return nil, err
	}

	return &watcher{
		lw:    w,
		rd:    c.opts.Reader,
		path:  path,
		value: value,
	}, nil
}

// Stop the config loader/watcher
func (c *config) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}

// Bytes get merged config data
func (c *config) Bytes() []byte {
	return c.vals.Bytes()
}

// Mount fields
func (c *config) Mount(fields ...*Field) {
	c.Lock()
	defer c.Unlock()

	c.userFields.Append(fields...).Sort()
}

func (c *config) MountSystem(fields ...*Field) {
	c.Lock()
	defer c.Unlock()

	c.sysFields.Append(fields...).Sort()
}

// Get value through field
func (c *config) Get(field *Field) reader.Value {
	c.RLock()
	defer c.RUnlock()

	// did sync actually work?
	if c.vals != nil {
		return c.vals.Get(field2path(field)...)
	}

	// no value
	return newValue()
}

// GetString through field
func (c *config) GetBool(field *Field) (val bool) {
	if field == nil {
		return
	}

	def, ok := field.def.(bool)
	if !ok {
		return
	}

	return c.Get(field).Bool(def)
}

// GetInt through field
func (c *config) GetInt(field *Field) (val int) {
	if field == nil {
		return
	}

	def, ok := field.def.(int)
	if !ok {
		return
	}

	return c.Get(field).Int(def)
}

// GetUint through field
func (c *config) GetUint(field *Field) (val uint) {
	if field == nil {
		return
	}

	def, ok := field.def.(uint)
	if !ok {
		return
	}

	return c.Get(field).Uint(def)
}

// GetString through field
func (c *config) GetString(field *Field) (val string) {
	if field == nil {
		return
	}

	def, ok := field.def.(string)
	if !ok {
		return
	}

	return c.Get(field).String(def)
}

// GetFloat64 through field
func (c *config) GetFloat64(field *Field) (val float64) {
	if field == nil {
		return
	}

	def, ok := field.def.(float64)
	if !ok {
		return
	}

	return c.Get(field).Float64(def)
}

// GetDuration through field
func (c *config) GetDuration(field *Field) (val time.Duration) {
	if field == nil {
		return
	}

	def, ok := field.def.(time.Duration)
	if !ok {
		return
	}

	return c.Get(field).Duration(def)
}

// GetStringSlice through field
func (c *config) GetStringSlice(field *Field) (val []string) {
	if field == nil {
		return
	}

	def, ok := field.def.([]string)
	if !ok {
		return
	}

	return c.Get(field).StringSlice(def)
}

// GetStringMap through field
func (c *config) GetStringMap(field *Field) (val map[string]string) {
	if field == nil {
		return
	}

	def, ok := field.def.(map[string]string)
	if !ok {
		return
	}

	return c.Get(field).StringMap(def)
}

// GetBoxName path: box.name
func (c *config) GetBoxName() string {
	return c.GetString(fieldBoxName)
}

// GetTraceUid path: box.trace.uid
func (c *config) GetTraceUid() string {
	return c.GetString(fieldTraceUid)
}

// GetTraceReqId path: box.trace.reqid
func (c *config) GetTraceReqId() string {
	return c.GetString(fieldTraceReqId)
}

// GetTraceBizId path: box.trace.bizid
func (c *config) GetTraceBizId() string {
	return c.GetString(fieldTraceBizId)
}

// GetTraceSpanId path: box.trace.spanid
func (c *config) GetTraceSpanId() string {
	return c.GetString(fieldTraceSpanId)
}

// SprintFields registered fields
func (c *config) SprintFields() (str string) {
	return sprintFields(c.sysFields, c.userFields)
}

// SprintTemplate through encoder
func (c *config) SprintTemplate(encoder string) (str string) {
	return sprintTemplate(c.sysFields, c.userFields, encoder)
}

func (c *config) String() string {
	return "config"
}
