package box

import (
	"context"
	"fmt"
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

	boxErr struct {
		box Box
		err error
	}
)

func (app *application) Run() error {
	if err := app.serve(); err != nil {
		return err
	}

	defer close(app.quit)
	<-app.quit

	if err := app.shutdown(); err != nil {
		return err
	}

	return nil
}

func (app *application) serve() error {
	logger.Info("serve start...")

	g := errgroup.Group{}
	for _, box := range app.boxes {
		box := box
		g.Go(func() error {
			serveCh := make(chan boxErr)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.startupTimeout)*time.Millisecond)
			defer cancel()

			go func() {
				defer close(serveCh)

				logger.Infof("serve %s", box.Name())
				serveCh <- boxErr{
					box: box,
					err: box.Serve(ctx),
				}
			}()

			select {
			case ch := <-serveCh:
				if ch.err != nil {
					logger.Errorf("serve %s error: %v", ch.box.Name(), ch.err)
				}

				return ch.err
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					return nil
				} else {
					return ctx.Err()
				}
			}
		})
	}

	err := g.Wait()
	if err != nil {
		logger.Errorf("serve error: %v", err)
	} else {
		logger.Infof("serve done...")
	}

	return err
}

func (app *application) shutdown() error {
	logger.Info("shutdown start...")

	g := errgroup.Group{}
	for _, box := range app.boxes {
		box := box
		g.Go(func() error {
			shutdownCh := make(chan boxErr)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.shutdownTimeout)*time.Millisecond)
			defer cancel()

			go func() {
				defer close(shutdownCh)

				logger.Infof("shutdown %s", box.Name())
				shutdownCh <- boxErr{
					box: box,
					err: box.Shutdown(ctx),
				}
			}()

			select {
			case ch := <-shutdownCh:
				if ch.err != nil {
					logger.Errorf("shutdown %s error: %v", ch.box.Name(), ch.err)
				}
				return ch.err
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}

	err := g.Wait()
	if err != nil {
		logger.Errorf("shutdown error: %v", err)
	} else {
		logger.Infof("shutdown done...")
	}

	return err
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
	if !opts.Silent {
		fmt.Print(banner)
		fmt.Print(config.SprintFields())
	}

	app := &application{
		quit:            make(chan os.Signal),
		startupTimeout:  opts.StartupTimeout,
		shutdownTimeout: opts.ShutdownTimeout,
		cfg:             opts.Configurator,
		boxes:           opts.Boxes,
	}

	if app.cfg == nil {
		logger.Fatal("configurator is required")
	} else if err := app.cfg.Sync(); err != nil {
		logger.Fatalf("configurator sync error: %v\n", err)
	} else {
		logger.Infof("configurator: %s", app.cfg.Bytes())
	}

	signal.Notify(app.quit, syscall.SIGINT, syscall.SIGTERM)

	return app
}
