package config

import (
	"bytes"

	"github.com/boxgo/box/pkg/config/loader"
	"github.com/boxgo/box/pkg/config/reader"
)

type (
	// Watcher is the config watcher
	Watcher interface {
		Next() (reader.Value, error)
		Stop() error
	}

	watcher struct {
		lw    loader.Watcher
		rd    reader.Reader
		path  []string
		value reader.Value
	}
)

func (w *watcher) Next() (reader.Value, error) {
	for {
		s, err := w.lw.Next()
		if err != nil {
			return nil, err
		}

		// only process changes
		if bytes.Equal(w.value.Bytes(), s.ChangeSet.Data) {
			continue
		}

		v, err := w.rd.Values(s.ChangeSet)
		if err != nil {
			return nil, err
		}

		w.value = v.Get()
		return w.value, nil
	}
}

func (w *watcher) Stop() error {
	return w.lw.Stop()
}
