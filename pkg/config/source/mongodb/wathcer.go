package mongodb

import (
	"log"
	"time"

	"github.com/boxgo/box/pkg/config/source"
)

type watcher struct {
	name      string
	source    *mongoSource
	exit      chan bool
	changeSet chan *source.ChangeSet
}

func newWatcher(mgo *mongoSource) (source.Watcher, error) {
	w := &watcher{
		name:      "mongodb",
		source:    mgo,
		changeSet: make(chan *source.ChangeSet),
		exit:      make(chan bool),
	}
	go w.watch()

	return w, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-w.changeSet:
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
				log.Printf("config mongodb watch error: %#v", err)
			}

			w.changeSet <- data
		}
	}
}
