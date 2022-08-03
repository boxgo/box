package box

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
	"github.com/boxgo/box/v2/server"
	"github.com/boxgo/box/v2/util/procsutil"
	"golang.org/x/sync/errgroup"
)

type (
	Application interface {
		Run() error
	}

	Box interface {
		Name() string
	}

	// application app interface
	application struct {
		startupTimeout  int
		shutdownTimeout int
		boxes           []Box
		quit            chan os.Signal
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
	for _, b := range app.boxes {
		b := b
		g.Go(func() error {
			var (
				hook        server.Hook
				ctx, cancel = context.WithTimeout(context.Background(), time.Duration(app.startupTimeout)*time.Millisecond)
			)
			defer cancel()

			if v, ok := b.(server.Hook); ok {
				hook = v
			}

			if hook != nil {
				if err := hook.BeforeServe(ctx); err != nil {
					return err
				}
			}

			if err := app.serveBox(ctx, b); err != nil {
				return err
			}

			if hook != nil {
				if err := hook.AfterServe(ctx); err != nil {
					return err
				}
			}

			return nil
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

func (app *application) serveBox(ctx context.Context, b Box) error {
	var (
		serv    server.Server
		serveCh = make(chan boxErr)
	)

	if v, ok := b.(server.Server); ok {
		serv = v
	}

	go func() {
		defer close(serveCh)

		if serv != nil {
			logger.Infof("serve %s", serv.Name())
			serveCh <- boxErr{
				box: b,
				err: serv.Serve(ctx),
			}
		}
	}()

	select {
	case ch := <-serveCh:
		if ch.err != nil {
			logger.Errorf("serve %s error: %v", ch.box.Name(), ch.err)
		}

		return ch.err
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			return ctx.Err()
		}
	}

	return nil
}

func (app *application) shutdown() error {
	logger.Info("shutdown start...")

	g := errgroup.Group{}
	for _, b := range app.boxes {
		b := b
		g.Go(func() error {
			var (
				hook        server.Hook
				ctx, cancel = context.WithTimeout(context.Background(), time.Duration(app.shutdownTimeout)*time.Millisecond)
			)
			defer cancel()

			if v, ok := b.(server.Hook); ok {
				hook = v
			}

			if hook != nil {
				if err := hook.BeforeShutdown(ctx); err != nil {
					return err
				}
			}

			if err := app.shutdownBox(ctx, b); err != nil {
				return err
			}

			if hook != nil {
				if err := hook.AfterShutdown(ctx); err != nil {
					return err
				}
			}

			return nil
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

func (app *application) shutdownBox(ctx context.Context, b Box) error {
	var (
		serv       server.Server
		shutdownCh = make(chan boxErr)
	)

	if v, ok := b.(server.Server); ok {
		serv = v
	}

	go func() {
		defer close(shutdownCh)

		if serv != nil {
			logger.Infof("shutdown %s", serv.Name())
			shutdownCh <- boxErr{
				box: b,
				err: serv.Shutdown(ctx),
			}
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

	if opts.AutoMaxProcs == nil || *opts.AutoMaxProcs {
		procsutil.EnableAutoMaxProcs()
	}

	config.AppendServiceTag(opts.Tags...)

	app := &application{
		quit:            make(chan os.Signal),
		startupTimeout:  opts.StartupTimeout,
		shutdownTimeout: opts.ShutdownTimeout,
		boxes:           append(opts.Boxes, &boxMetric{}),
	}

	signal.Notify(app.quit, syscall.SIGINT, syscall.SIGTERM)

	return app
}
