package box

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/boxgo/box/config"
	"github.com/boxgo/box/minibox"
	"github.com/boxgo/kit/logger"
)

type (
	// Box is a application bootstrap
	Box struct {
		App string `json:"name" desc:"Application name"`

		ctx        context.Context
		boxes      []minibox.MiniBox
		config     config.Config
		configHook minibox.ConfigHook
		serverHook minibox.ServerHook
	}

	// Options new box options
	Options struct {
		Context    context.Context
		Boxes      []minibox.MiniBox
		Config     config.Config
		ConfigHook minibox.ConfigHook
		ServerHook minibox.ServerHook
	}

	// Option setter
	Option func(ops *Options)
)

// Name Box is a minibox too.
func (box *Box) Name() string {
	return ""
}

// Mount mount minibox to box
func (box *Box) Mount(boxes ...minibox.MiniBox) *Box {
	box.boxes = append(box.boxes, boxes...)

	return box
}

// Serve all miniboxes
func (box *Box) Serve() error {
	box.setupConfig()

	if box.serverHook != nil {
		box.serverHook.ServerWillReady(box.ctx)
	}

	box.traverseMiniBoxes(func(mini minibox.MiniBox) {
		var serverOk, serverHookOk bool
		var server minibox.Server
		var serverHook minibox.ServerHook

		server, serverOk = mini.(minibox.Server)
		serverHook, serverHookOk = mini.(minibox.ServerHook)

		if serverHookOk {
			serverHook.ServerWillReady(box.ctx)
		}

		if serverOk {
			if err := server.Serve(box.ctx); err != nil {
				logger.Default.Errorf("MiniBox [%s] serve error: %s", mini.Name(), err)
				panic(err)
			} else {
				logger.Default.Infof("MiniBox [%s] serve success", mini.Name())
			}
		}

		if serverHookOk {
			serverHook.ServerDidReady(box.ctx)
		}
	})

	if box.serverHook != nil {
		box.serverHook.ServerDidReady(box.ctx)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	box.Shutdown()

	return nil
}

// Shutdown all miniboxes
func (box *Box) Shutdown() {
	if box.serverHook != nil {
		box.serverHook.ServerWillClose(box.ctx)
	}

	box.traverseMiniBoxes(func(mini minibox.MiniBox) {
		var serverOk, serverHookOk bool
		var server minibox.Server
		var serverHook minibox.ServerHook

		server, serverOk = mini.(minibox.Server)
		serverHook, serverHookOk = mini.(minibox.ServerHook)

		if serverHookOk {
			serverHook.ServerWillClose(box.ctx)
		}

		if serverOk {
			if err := server.Serve(box.ctx); err != nil {
				logger.Default.Errorf("MiniBox [%s] shutdown error: %s", mini.Name(), err)
				panic(err)
			} else {
				logger.Default.Infof("MiniBox [%s] shutdown success", mini.Name())
			}
		}

		if serverHookOk {
			serverHook.ServerDidClose(box.ctx)
		}
	})

	if box.serverHook != nil {
		box.serverHook.ServerDidClose(box.ctx)
	}
}

func (box *Box) traverseMiniBoxes(done func(minibox.MiniBox)) {
	wg := sync.WaitGroup{}

	for _, b := range box.boxes {
		wg.Add(1)

		go func(miniBox minibox.MiniBox) {
			done(miniBox)
			wg.Done()
		}(b)
	}

	wg.Wait()
}

func (box *Box) setupConfig() {
	if box.configHook != nil {
		box.configHook.ConfigWillLoad(box.ctx)
	}

	box.config.Set(box.boxes...).Scan(box.ctx)

	if box.configHook != nil {
		box.configHook.ConfigDidLoad(box.ctx)
	}
}

// NewBox new a box bootstrap
func NewBox(options ...Option) *Box {
	opts := &Options{}

	for _, opt := range options {
		opt(opts)
	}

	if opts.Config == nil {
		panic("Config must be set.")
	}

	b := &Box{
		ctx:        context.Background(),
		config:     opts.Config,
		configHook: opts.ConfigHook,
		serverHook: opts.ServerHook,
	}

	b.Mount(b).Mount(opts.Boxes...)

	return b
}
