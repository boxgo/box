package gopg

import (
	"context"
	"io"
	"time"

	"github.com/boxgo/box/pkg/logger"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type (
	PostgreSQL struct {
		db *pg.DB
	}
)

func newPostgreSQL(c *Config) *PostgreSQL {
	var (
		err  error
		opts *pg.Options
	)

	if c.URI != "" {
		opts, err = pg.ParseURL(c.URI)
		if err != nil {
			logger.Panicf("PostgreSQL ParseURL [%s] error: %s", c.URI, err)
		}
	} else {
		opts = &pg.Options{
			Network:               c.Network,
			Addr:                  c.Addr,
			Dialer:                c.dialer,
			OnConnect:             c.onConnect,
			User:                  c.User,
			Password:              c.Password,
			Database:              c.Database,
			ApplicationName:       c.ApplicationName,
			TLSConfig:             c.tLSConfig,
			DialTimeout:           c.DialTimeout,
			ReadTimeout:           c.ReadTimeout,
			WriteTimeout:          c.WriteTimeout,
			MaxRetries:            c.MaxRetries,
			RetryStatementTimeout: c.RetryStatementTimeout,
			MinRetryBackoff:       c.MinRetryBackoff,
			MaxRetryBackoff:       c.MaxRetryBackoff,
			PoolSize:              c.PoolSize,
			MinIdleConns:          c.MinIdleConns,
			MaxConnAge:            c.MaxConnAge,
			PoolTimeout:           c.PoolTimeout,
			IdleTimeout:           c.IdleTimeout,
			IdleCheckFrequency:    c.IdleCheckFrequency,
		}
	}

	db := pg.Connect(opts)

	if c.Debug {
		db.AddQueryHook(DebugHook{Verbose: true})
	}

	return &PostgreSQL{
		db: db,
	}
}

func (pg *PostgreSQL) Name() string {
	return "gopg"
}

func (pg *PostgreSQL) Serve(ctx context.Context) error {
	if err := pg.db.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func (pg *PostgreSQL) Shutdown(ctx context.Context) error {
	if pg.db != nil {
		return pg.db.Close()
	}

	return nil
}

func (pg *PostgreSQL) Begin() (*pg.Tx, error) {
	return pg.db.Begin()
}

func (pg *PostgreSQL) BeginContext(ctx context.Context) (*pg.Tx, error) {
	return pg.db.BeginContext(ctx)
}

func (pg *PostgreSQL) Conn() *pg.Conn {
	return pg.db.Conn()
}

func (pg *PostgreSQL) Context() context.Context {
	return pg.db.Context()
}

func (pg *PostgreSQL) CopyFrom(r io.Reader, query interface{}, params ...interface{}) (res pg.Result, err error) {
	return pg.db.CopyFrom(r, query, params...)
}

func (pg *PostgreSQL) CopyTo(w io.Writer, query interface{}, params ...interface{}) (res pg.Result, err error) {
	return pg.db.CopyTo(w, query, params...)
}

func (pg *PostgreSQL) Exec(query interface{}, params ...interface{}) (res pg.Result, err error) {
	return pg.db.Exec(query, params...)
}

func (pg *PostgreSQL) ExecContext(c context.Context, query interface{}, params ...interface{}) (pg.Result, error) {
	return pg.db.ExecContext(c, query, params)
}

func (pg *PostgreSQL) ExecOne(query interface{}, params ...interface{}) (pg.Result, error) {
	return pg.db.ExecOne(query, params...)
}

func (pg *PostgreSQL) ExecOneContext(ctx context.Context, query interface{}, params ...interface{}) (pg.Result, error) {
	return pg.db.ExecOneContext(ctx, query, params...)
}

func (pg *PostgreSQL) Formatter() orm.QueryFormatter {
	return pg.db.Formatter()
}

func (pg *PostgreSQL) Listen(ctx context.Context, channels ...string) *pg.Listener {
	return pg.db.Listen(ctx, channels...)
}

func (pg *PostgreSQL) Model(model ...interface{}) *orm.Query {
	return pg.db.Model(model...)
}

func (pg *PostgreSQL) ModelContext(c context.Context, model ...interface{}) *orm.Query {
	return pg.db.ModelContext(c, model...)
}

func (pg *PostgreSQL) Param(param string) interface{} {
	return pg.db.Param(param)
}

func (pg *PostgreSQL) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *PostgreSQL) PoolStats() *pg.PoolStats {
	return pg.db.PoolStats()
}

func (pg *PostgreSQL) Prepare(q string) (*pg.Stmt, error) {
	return pg.db.Prepare(q)
}

func (pg *PostgreSQL) Query(model, query interface{}, params ...interface{}) (res pg.Result, err error) {
	return pg.db.Query(model, query, params)
}

func (pg *PostgreSQL) QueryContext(c context.Context, model, query interface{}, params ...interface{}) (pg.Result, error) {
	return pg.db.QueryContext(c, model, query, params...)
}

func (pg *PostgreSQL) QueryOne(model, query interface{}, params ...interface{}) (pg.Result, error) {
	return pg.db.QueryOne(model, query, params...)
}

func (pg *PostgreSQL) QueryOneContext(ctx context.Context, model, query interface{}, params ...interface{}) (pg.Result, error) {
	return pg.db.QueryOneContext(ctx, model, query, params...)
}

func (pg *PostgreSQL) RunInTransaction(ctx context.Context, fn func(*pg.Tx) error) error {
	return pg.db.RunInTransaction(ctx, fn)
}

func (pg *PostgreSQL) WithContext(ctx context.Context) *pg.DB {
	return pg.db.WithContext(ctx)
}

func (pg *PostgreSQL) WithParam(param string, value interface{}) *pg.DB {
	return pg.db.WithParam(param, value)
}

func (pg *PostgreSQL) WithTimeout(d time.Duration) *pg.DB {
	return pg.db.WithTimeout(d)
}
