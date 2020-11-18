package mongodb

import (
	"bytes"
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/boxgo/box/pkg/config/source"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type watcher struct {
	sync.RWMutex
	name       string
	opts       source.Options
	cs         *source.ChangeSet
	rsp        chan string
	ch         chan *source.ChangeSet
	exit       chan bool
	client     *mongo.Client
	db         string
	collection string
}

func newWatcher(db, collection string, client *mongo.Client, opts source.Options) (source.Watcher, error) {
	w := &watcher{
		name:       "mongo",
		db:         db,
		collection: collection,
		opts:       opts,
		cs:         nil,
		rsp:        make(chan string),
		ch:         make(chan *source.ChangeSet),
		exit:       make(chan bool),
		client:     client,
	}
	go w.watch()

	return w, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-w.ch:
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

			cfg := &Config{}
			err := w.client.Database(w.db).
				Collection(w.collection).
				FindOne(ctx, bson.D{}).
				Decode(cfg)

			if err != nil {
				log.Printf("config redis watch error: %#v", err)
			} else if cfg.Config != "" {
				w.handle(cfg)
			}

			cancel()
		}
	}
}

func (w *watcher) handle(cfg *Config) {
	w.RLock()
	eq := w.cs != nil && bytes.Compare(w.cs.Data, []byte(cfg.Config)) == 0
	w.RUnlock()

	if eq {
		return
	}

	var val map[string]interface{}
	if err := w.opts.Encoder.Decode([]byte(cfg.Config), &val); err != nil {
		log.Printf("config mongo watch handler decode error: %#v", err)
		return
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Source:    w.name,
		Data:      []byte(cfg.Config),
		Format:    w.opts.Encoder.String(),
	}
	cs.Checksum = cs.Sum()

	w.Lock()
	w.cs = cs
	w.Unlock()

	w.ch <- cs
}
