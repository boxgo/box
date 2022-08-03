package redis

import (
	"errors"
	"log"
	"time"

	"github.com/boxgo/box/v2/config/source"
)

type watcher struct {
	name       string
	source     *redisSource
	exit       chan bool
	changeSets chan *source.ChangeSet
}

func newWatcher(sour *redisSource) (source.Watcher, error) {
	w := &watcher{
		name:       "redis",
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
		return nil, errors.New("watcher stopped")
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
				log.Printf("config redis watch error: %#v", err)
			}

			w.changeSets <- data
		}
	}
}
