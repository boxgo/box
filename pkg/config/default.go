package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/boxgo/box/pkg/app"
	"github.com/boxgo/box/pkg/config/field"
	"github.com/boxgo/box/pkg/config/loader"
	"github.com/boxgo/box/pkg/config/loader/memory"
	"github.com/boxgo/box/pkg/config/reader"
	"github.com/boxgo/box/pkg/config/reader/json"
	"github.com/boxgo/box/pkg/config/source"
	"github.com/boxgo/box/pkg/config/util"
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
		sysFields field.Fields
		// user registered fields
		userFields field.Fields
	}
)

var (
	traceUid    = field.New(true, "trace", "uid", "trace uid in context", "box.trace.uid")
	traceReqId  = field.New(true, "trace", "reqId", "trace requestId in context", "box.trace.reqId")
	traceSpanId = field.New(true, "trace", "spanId", "trace spanId in context", "box.trace.spanId")
	traceBizId  = field.New(true, "trace", "bizId", "trace bizId in context", "box.trace.bizId")
)

func newConfig(opts ...Option) Configurator {
	options := Options{
		Loader: memory.NewLoader(),
		Reader: json.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	if err := options.Loader.Load(options.Source...); err != nil {
		panic(err)
	}

	snap, _ := options.Loader.Snapshot()
	vals, _ := options.Reader.Values(snap.ChangeSet)

	c := &config{
		exit:       make(chan bool),
		opts:       options,
		snap:       snap,
		vals:       vals,
		sysFields:  make(field.Fields, 0),
		userFields: make(field.Fields, 0),
	}

	c.MountSystem(
		traceUid,
		traceReqId,
		traceSpanId,
		traceBizId,
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
func (c *config) Watch(field *field.Field) (Watcher, error) {
	value := c.Get(field)

	if field.Immutable {
		return nil, fmt.Errorf("field [%s] is immutable", field.String())
	}

	w, err := c.opts.Loader.Watch(field.Paths()...)
	if err != nil {
		return nil, err
	}

	return &watcher{
		lw:    w,
		rd:    c.opts.Reader,
		path:  field.Paths(),
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
func (c *config) Mount(fields ...*field.Field) {
	c.Lock()
	defer c.Unlock()

	c.userFields.Append(fields...).Sort()
}

func (c *config) MountSystem(fields ...*field.Field) {
	c.Lock()
	defer c.Unlock()

	c.sysFields.Append(fields...).Sort()
}

// Get value through field
func (c *config) Get(field *field.Field) reader.Value {
	c.RLock()
	defer c.RUnlock()

	if field.Immutable && field.Val() != nil {
		return field.Val()
	}

	// did sync actually work?
	if c.vals != nil {
		field.SetVal(c.vals.Get(field.Paths()...))

		return field.Val()
	}

	// no value
	return newValue()
}

// GetString through field
func (c *config) GetBool(field *field.Field) (val bool) {
	if field == nil {
		return
	}

	def, ok := field.Def.(bool)
	if !ok {
		return
	}

	return c.Get(field).Bool(def)
}

// GetInt through field
func (c *config) GetInt(field *field.Field) (val int) {
	if field == nil {
		return
	}

	def, ok := field.Def.(int)
	if !ok {
		return
	}

	return c.Get(field).Int(def)
}

// GetUint through field
func (c *config) GetUint(field *field.Field) (val uint) {
	if field == nil {
		return
	}

	def, ok := field.Def.(uint)
	if !ok {
		return
	}

	return c.Get(field).Uint(def)
}

// GetString through field
func (c *config) GetString(field *field.Field) (val string) {
	if field == nil {
		return
	}

	def, ok := field.Def.(string)
	if !ok {
		return
	}

	return c.Get(field).String(def)
}

// GetFloat64 through field
func (c *config) GetFloat64(field *field.Field) (val float64) {
	if field == nil {
		return
	}

	def, ok := field.Def.(float64)
	if !ok {
		return
	}

	return c.Get(field).Float64(def)
}

// GetDuration through field
func (c *config) GetDuration(field *field.Field) (val time.Duration) {
	if field == nil {
		return
	}

	def, ok := field.Def.(time.Duration)
	if !ok {
		return
	}

	return c.Get(field).Duration(def)
}

// GetStringSlice through field
func (c *config) GetStringSlice(field *field.Field) (val []string) {
	if field == nil {
		return
	}

	def, ok := field.Def.([]string)
	if !ok {
		return
	}

	return c.Get(field).StringSlice(def)
}

// GetStringMap through field
func (c *config) GetStringMap(field *field.Field) (val map[string]string) {
	if field == nil {
		return
	}

	def, ok := field.Def.(map[string]string)
	if !ok {
		return
	}

	return c.Get(field).StringMap(def)
}

// GetBoxName path: box.name
func (c *config) GetBoxName() string {
	return app.Name
}

// GetTraceUid path: box.trace.uid
func (c *config) GetTraceUid() string {
	return c.GetString(traceUid)
}

// GetTraceReqId path: box.trace.reqid
func (c *config) GetTraceReqId() string {
	return c.GetString(traceReqId)
}

// GetTraceBizId path: box.trace.bizid
func (c *config) GetTraceBizId() string {
	return c.GetString(traceBizId)
}

// GetTraceSpanId path: box.trace.spanid
func (c *config) GetTraceSpanId() string {
	return c.GetString(traceSpanId)
}

// SprintFields registered fields
func (c *config) SprintFields() (str string) {
	return util.SprintFields(c.sysFields, c.userFields)
}

func (c *config) String() string {
	return "config"
}
