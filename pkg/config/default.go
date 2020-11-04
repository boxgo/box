package config

import (
	"sync"
	"time"

	"github.com/boxgo/box/pkg/config/field"
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
		// scanned fields
		fields field.Fields
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

	if err := options.Loader.Load(options.Source...); err != nil {
		panic(err)
	}

	snap, _ := options.Loader.Snapshot()
	vals, _ := options.Reader.Values(snap.ChangeSet)

	c := &config{
		exit: make(chan bool),
		opts: options,
		snap: snap,
		vals: vals,
	}

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

// Scan config to val
func (c *config) Scan(val Config) error {
	c.RLock()
	defer c.RUnlock()

	c.fields.Parse(val)

	return c.vals.Get(val.Path()).Scan(val)
}

// Watch a value for changes
func (c *config) Watch(path ...string) (Watcher, error) {
	value := c.Get(path...)

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

// Get value through field
func (c *config) Get(path ...string) reader.Value {
	c.RLock()
	defer c.RUnlock()

	// did sync actually work?
	if c.vals != nil {
		return c.vals.Get(path...)
	}

	// no value
	return newValue()
}

func (c *config) Fields() *field.Fields {
	return &c.fields
}

func (c *config) String() string {
	return "config"
}
