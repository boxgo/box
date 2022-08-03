package pubsub

import (
	"errors"
	"sync"
)

type (
	Message interface{}

	MessageChan chan Message

	MessageCallback func(Message)

	Options struct {
		bufferSize int
	}

	OptionFunc func(*Options)

	PubSub struct {
		close    bool
		options  *Options
		waitGp   sync.WaitGroup
		rwMutex  sync.RWMutex
		handlers map[string][]handler
	}

	handler struct {
		ch MessageChan
		cb MessageCallback
	}
)

var (
	ErrClosed      = errors.New("pub/sub is closed")
	ErrNoSubscribe = errors.New("no subscribe")
)

func New(optionFunc ...OptionFunc) *PubSub {
	opts := &Options{
		bufferSize: 1000,
	}
	for _, fn := range optionFunc {
		fn(opts)
	}

	return &PubSub{
		close:    false,
		options:  opts,
		waitGp:   sync.WaitGroup{},
		rwMutex:  sync.RWMutex{},
		handlers: make(map[string][]handler),
	}
}

func WithBufferSize(size int) OptionFunc {
	return func(options *Options) {
		options.bufferSize = size
	}
}

func (pb *PubSub) Subscribe(topic string, callback MessageCallback) error {
	if pb.close {
		return ErrClosed
	}

	hdr := pb.register(topic, callback)

	go func() {
		for data := range hdr.ch {
			hdr.cb(data)
			pb.waitGp.Done()
		}
	}()

	return nil
}

func (pb *PubSub) Publish(topic string, data interface{}) error {
	pb.rwMutex.RLock()
	defer pb.rwMutex.RUnlock()

	if pb.close {
		return ErrClosed
	}

	handlers, ok := pb.handlers[topic]
	if !ok {
		return ErrNoSubscribe
	}

	for _, hdr := range handlers {
		pb.waitGp.Add(1)
		hdr.ch <- data
	}

	return nil
}

func (pb *PubSub) Close() {
	pb.rwMutex.Lock()
	defer pb.rwMutex.Unlock()

	pb.close = true

	for _, handlers := range pb.handlers {
		for _, h := range handlers {
			close(h.ch)
		}
	}
}

func (pb *PubSub) Wait() {
	pb.waitGp.Wait()
}

func (pb *PubSub) register(topic string, callback MessageCallback) *handler {
	pb.rwMutex.Lock()
	defer pb.rwMutex.Unlock()

	_, ok := pb.handlers[topic]
	h := handler{
		ch: make(MessageChan, pb.options.bufferSize),
		cb: callback,
	}

	if !ok {
		pb.handlers[topic] = []handler{h}
	} else {
		pb.handlers[topic] = append(pb.handlers[topic], h)
	}

	return &h
}
