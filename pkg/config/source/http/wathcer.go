package http

import (
	"log"
	"time"

	"github.com/boxgo/box/pkg/config/source"
)

type watcher struct {
	name       string
	source     *httpSource
	exit       chan bool
	changeSets chan *source.ChangeSet
}

func newWatcher(sour *httpSource) (source.Watcher, error) {
	w := &watcher{
		name:       "http",
		changeSets: make(chan *source.ChangeSet),
		exit:       make(chan bool),
		source:     sour,
	}
	go w.watch()

	return w, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-w.changeSets:
		return cs, nil
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
		return nil
	default:
		close(w.exit)
	}
	return nil
}

func (w *watcher) watch() {
	for {
		time.Sleep(time.Second * 3)

		select {
		case <-w.exit:
			return
		default:
			data, err := w.source.Read()
			if err != nil {
				log.Printf("config http watch error: %#v", err)
			}

			w.changeSets <- data
		}
	}
}
