package redis

import (
	"bytes"
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/boxgo/box/pkg/config/source"
	"github.com/go-redis/redis/v8"
)

type watcher struct {
	sync.RWMutex
	name   string
	prefix string
	opts   source.Options
	cs     *source.ChangeSet
	rsp    chan string
	ch     chan *source.ChangeSet
	exit   chan bool
	client redis.UniversalClient
}

func newWatcher(prefix string, client redis.UniversalClient, opts source.Options) (source.Watcher, error) {
	w := &watcher{
		name:   "redis",
		prefix: prefix,
		opts:   opts,
		cs:     nil,
		rsp:    make(chan string),
		ch:     make(chan *source.ChangeSet),
		exit:   make(chan bool),
		client: client,
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
			if rsp, err := w.client.Get(ctx, w.prefix+".config").Bytes(); err != nil && err != redis.Nil {
				log.Printf("config redis watch error: %#v", err)
			} else if len(rsp) != 0 {
				w.handle(rsp)
			}

			cancel()
		}
	}
}

func (w *watcher) handle(data []byte) {
	w.RLock()
	eq := w.cs != nil && bytes.Compare(w.cs.Data, data) == 0
	w.RUnlock()

	if eq {
		return
	}

	var val map[string]interface{}
	if err := w.opts.Encoder.Decode(data, &val); err != nil {
		log.Printf("config redis watch handler decode error: %#v", err)
		return
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Source:    w.name,
		Data:      data,
		Format:    w.opts.Encoder.String(),
	}
	cs.Checksum = cs.Sum()

	w.Lock()
	w.cs = cs
	w.Unlock()

	w.ch <- cs
}
