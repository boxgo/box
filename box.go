package box

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/server"
	"golang.org/x/sync/errgroup"
)

type (
	// Box component interface
	Box interface {
		config.Config
		server.Server
	}

	Application interface {
		Run() error
	}

	// application app interface
	application struct {
		startupTimeout  int
		shutdownTimeout int
		boxes           []Box
		quit            chan os.Signal
		cfg             config.Configurator
	}
)

func (app *application) Run() error {
	logger.Info("box application run")
	if err := app.init(); err != nil {
		return err
	}

	if err := app.serve(); err != nil {
		return err
	}

	<-app.quit

	if err := app.shutdown(); err != nil {
		return err
	}

	return nil
}

func (app *application) init() error {
	logger.Debug("register signal")
	signal.Notify(app.quit, syscall.SIGINT, syscall.SIGTERM)

	// force sync config source
	if err := app.cfg.Sync(); err != nil {
		return err
	}

	logger.Infof("configs: %s", app.cfg.Bytes())

	// register config mounter getter to boxes
	for _, box := range app.boxes {
		logger.Infof("configure [%s]", box.Name())
		box.Configure(app.cfg)
	}

	return nil
}

func (app *application) serve() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.startupTimeout)*time.Millisecond)
	defer cancel()

	logger.Info("serve start...")

	g := errgroup.Group{}
	for _, box := range app.boxes {
		box := box
		g.Go(func() error {
			err := box.Serve(ctx)
			if err != nil {
				logger.Infof("serve [%s] error: %v", box.Name(), err)
			} else {
				logger.Infof("serve [%s] success", box.Name())
			}

			return err
		})
	}

	errCh := make(chan error)
	go func() {
		errCh <- g.Wait()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			logger.Errorf("serve error: %v", err)
		} else {
			logger.Info("serve success")
		}
		return err
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			logger.Infof("serve done...")
			return nil
		} else {
			logger.Errorf("serve server error: %v", ctx.Err())
			return ctx.Err()
		}
	}
}

func (app *application) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.shutdownTimeout)*time.Millisecond)
	defer cancel()

	logger.Info("shutdown start...")

	g := errgroup.Group{}
	for _, box := range app.boxes {
		box := box
		g.Go(func() error {
			err := box.Shutdown(ctx)
			if err != nil {
				logger.Infof("shutdown [%s] error: %v", box.Name(), err)
			} else {
				logger.Infof("shutdown [%s] success", box.Name())
			}

			return err
		})
	}

	errCh := make(chan error)
	go func() {
		errCh <- g.Wait()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			logger.Errorf("shutdown error: %v", err)
		} else {
			logger.Info("shutdown done...")
		}
		return err
	case <-ctx.Done():
		logger.Infof("shutdown timeout: %v", ctx.Err())
		return ctx.Err()
	}
}

// New new a box bootstrap
func New(options ...Option) Application {
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}

	if opts.StartupTimeout == 0 {
		opts.StartupTimeout = 1000
	}
	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 5000
	}

	app := &application{
		quit:            make(chan os.Signal),
		startupTimeout:  opts.StartupTimeout,
		shutdownTimeout: opts.ShutdownTimeout,
		cfg:             opts.Configurator,
		boxes:           opts.Boxes,
	}

	return app
}
