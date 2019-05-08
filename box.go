package box

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boxgo/box/config"
	"github.com/boxgo/box/minibox"
	"github.com/boxgo/logger"
)

type (
	// Box is a application bootstrap
	Box struct {
		minibox.App

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

	go func() {
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
				go func() {
					logger.Default.Infof("server %-20s ready to serve", mini.Name())
					if err := server.Serve(box.ctx); err != nil {
						logger.Default.Errorf("server %-20s serve error: %s", mini.Name(), err)
						panic(err)
					}
				}()
			}

			// 等待100毫秒启动，如果没有报错，认为是启动成功了
			time.Sleep(time.Millisecond * 100)
			if serverHookOk {
				serverHook.ServerDidReady(box.ctx)
			}
		})

		if box.serverHook != nil {
			box.serverHook.ServerDidReady(box.ctx)
		}
	}()

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
			if err := server.Shutdown(box.ctx); err != nil {
				logger.Default.Errorf("server %-20s shutdown error: %s", mini.Name(), err)
				panic(err)
			} else {
				logger.Default.Infof("server %-20s shutdown success", mini.Name())
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

func (box *Box) traverseMiniBoxes(handler func(minibox.MiniBox)) {
	for _, b := range box.boxes {
		handler(b)
	}
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
