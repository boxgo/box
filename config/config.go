// Package config is an interface for dynamic configuration.
package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/boxgo/box/v2/config/reader"
	"github.com/boxgo/box/v2/config/source"
	"github.com/boxgo/box/v2/config/value"
	"github.com/boxgo/box/v2/config/value/json"
)

type (
	// Configurator is an interface abstraction for dynamic configuration
	Configurator interface {
		Load() error                         // Load config sources
		Scan(val interface{}) error          // Scan to val
		Value(key string) value.Value        // Value get values through key
		Watch(key string, cb Callback) error // Watch key change
		Close() error                        // Close stop the config loader/watcher
	}

	Config interface {
		Path() string
	}

	Callback func(old, new value.Value)

	config struct {
		sync.RWMutex
		stop        chan bool
		opts        Options
		values      value.Values
		watchers    sync.Map
		watcherCBs  sync.Map
		watcherVals sync.Map
	}
)

// NewConfig returns new config
func NewConfig(opts ...Option) Configurator {
	options := Options{
		Debug:    io.Discard,
		Interval: time.Second,
		Reader:   reader.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	values, _ := json.NewValues(nil)

	c := &config{
		stop:        make(chan bool),
		opts:        options,
		values:      values,
		watchers:    sync.Map{},
		watcherCBs:  sync.Map{},
		watcherVals: sync.Map{},
	}

	c.run()

	return c
}

func (c *config) Load() error {
	c.Lock()
	defer c.Unlock()

	changes := make([]*source.ChangeSet, len(c.opts.Source))
	for idx, sour := range c.opts.Source {
		data, err := sour.Read()
		if err != nil {
			return err
		}

		changes[idx] = data
	}

	if mergeData, err := c.opts.Reader.Merge(changes...); err != nil {
		return err
	} else if c.values, err = json.NewValues(mergeData); err != nil {
		return err
	}

	return nil
}

func (c *config) Scan(val interface{}) (err error) {
	c.RLock()
	defer c.RUnlock()

	switch v := val.(type) {
	case Config:
		err = c.Value(v.Path()).Scan(val)
	default:
		err = c.values.Scan(val)
	}

	if err != nil {
		return
	}

	if c.opts.Validator != nil {
		err = c.opts.Validator.Validate(val)
	}

	return
}

func (c *config) Value(key string) value.Value {
	c.RLock()
	defer c.RUnlock()

	return c.values.Value(key)
}

func (c *config) Watch(key string, cb Callback) error {
	if val, ok := c.watcherCBs.Load(key); !ok {
		c.watcherCBs.Store(key, []Callback{cb})
	} else if watchers, ok := val.([]Callback); ok {
		c.watcherCBs.Store(key, append(watchers, cb))
	}

	c.watcherVals.Store(key, c.Value(key))

	return nil
}

func (c *config) Close() error {
	c.watchers.Range(func(key, value interface{}) bool {
		if w, ok := value.(source.Watcher); ok {
			fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Watch.Stop\n", key)
			w.Stop()
		}

		return true
	})

	c.stop <- true

	return nil
}

func (c *config) run() {
	for _, s := range c.opts.Source {
		s := s

		if w, err := s.Watch(); err != nil {
			fmt.Fprintf(os.Stderr, "Config.Source[%s].Watch.Error error:%s", s.Id(), err)
			return
		} else {
			c.watchers.Store(s.Id(), w)
			go func() {
				fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Watch.Start\n", s.Id())

			watchLoop:
				for {
					select {
					case <-c.stop:
						fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Watcher.Stop\n", s.Id())
						break watchLoop
					default:
						c.watch(w, s)
					}

					time.Sleep(c.opts.Interval)
				}

				defer w.Stop()
			}()
		}
	}
}

func (c *config) watch(w source.Watcher, s source.Source) {
	var (
		err error
	)

	fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Next.Start\n", s.Id())
	if _, err = w.Next(); err != nil {
		fmt.Fprintf(os.Stderr, "Config.Source[%s].Next.Error error:%s", s.Id(), err)
		return
	}

	fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Load.Start\n", s.Id())
	if err = c.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Config.Source[%s].Load.Error error:%s", s.Id(), err)
	}

	c.watcherCBs.Range(func(k, v interface{}) bool {
		key := k.(string)
		cbs := v.([]Callback)

		var (
			newValue = c.Value(key)
			oldValue value.Value
		)

		if oldInterface, oldOk := c.watcherVals.Load(key); oldOk {
			if val, ok := oldInterface.(value.Value); ok {
				oldValue = val
			}
		}

		if oldValue != nil && newValue != nil && bytes.Equal(oldValue.Bytes(), newValue.Bytes()) {
			fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Key[%s] skip same value\n", s.Id(), key)
			return true
		}

		if oldValue == nil {
			fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Key[%s] %s -> %s\n", s.Id(), key, "nil", newValue)
		} else {
			fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Key[%s] %s -> %s\n", s.Id(), key, oldValue, newValue)
		}

		for _, cb := range cbs {
			c.watcherVals.Store(key, newValue)
			cb(oldValue, newValue)
		}

		return true
	})

	fmt.Fprintf(c.opts.Debug, "Config.Source[%s].Read.Finish\n", s.Id())
}
