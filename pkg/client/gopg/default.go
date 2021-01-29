package gopg

import (
	"context"
	"io"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var (
	// Default instance
	Default = StdConfig("default").Build()
)

func Begin() (*pg.Tx, error) {
	return Default.Begin()
}

func BeginContext(ctx context.Context) (*pg.Tx, error) {
	return Default.BeginContext(ctx)
}

func Conn() *pg.Conn {
	return Default.Conn()
}

func Context() context.Context {
	return Default.Context()
}

func CopyFrom(r io.Reader, query interface{}, params ...interface{}) (res pg.Result, err error) {
	return Default.CopyFrom(r, query, params...)
}

func CopyTo(w io.Writer, query interface{}, params ...interface{}) (res pg.Result, err error) {
	return Default.CopyTo(w, query, params...)
}

func Exec(query interface{}, params ...interface{}) (res pg.Result, err error) {
	return Default.Exec(query, params...)
}

func ExecContext(c context.Context, query interface{}, params ...interface{}) (pg.Result, error) {
	return Default.ExecContext(c, query, params)
}

func ExecOne(query interface{}, params ...interface{}) (pg.Result, error) {
	return Default.ExecOne(query, params...)
}

func ExecOneContext(ctx context.Context, query interface{}, params ...interface{}) (pg.Result, error) {
	return Default.ExecOneContext(ctx, query, params...)
}

func Formatter() orm.QueryFormatter {
	return Default.Formatter()
}

func Listen(ctx context.Context, channels ...string) *pg.Listener {
	return Default.Listen(ctx, channels...)
}

func Model(model ...interface{}) *orm.Query {
	return Default.Model(model...)
}

func ModelContext(c context.Context, model ...interface{}) *orm.Query {
	return Default.ModelContext(c, model...)
}

func Param(param string) interface{} {
	return Default.Param(param)
}

func Ping(ctx context.Context) error {
	return Default.Ping(ctx)
}

func PoolStats() *pg.PoolStats {
	return Default.PoolStats()
}

func Prepare(q string) (*pg.Stmt, error) {
	return Default.Prepare(q)
}

func Query(model, query interface{}, params ...interface{}) (res pg.Result, err error) {
	return Default.Query(model, query, params)
}

func QueryContext(c context.Context, model, query interface{}, params ...interface{}) (pg.Result, error) {
	return Default.QueryContext(c, model, query, params...)
}

func QueryOne(model, query interface{}, params ...interface{}) (pg.Result, error) {
	return Default.QueryOne(model, query, params...)
}

func QueryOneContext(ctx context.Context, model, query interface{}, params ...interface{}) (pg.Result, error) {
	return Default.QueryOneContext(ctx, model, query, params...)
}

func RunInTransaction(ctx context.Context, fn func(*pg.Tx) error) error {
	return Default.RunInTransaction(ctx, fn)
}

func WithContext(ctx context.Context) *pg.DB {
	return Default.WithContext(ctx)
}

func WithParam(param string, value interface{}) *pg.DB {
	return Default.WithParam(param, value)
}

func WithTimeout(d time.Duration) *pg.DB {
	return Default.WithTimeout(d)
}
