package gopg

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type (
	Error                  = pg.Error
	DB                     = pg.DB
	PgConn                 = pg.Conn
	Tx                     = pg.Tx
	Result                 = pg.Result
	Stmt                   = pg.Stmt
	Options                = pg.Options
	Ident                  = pg.Ident
	Safe                   = pg.Safe
	NullTime               = pg.NullTime
	PgQuery                = orm.Query
	CreateTableOptions     = orm.CreateTableOptions
	CreateCompositeOptions = orm.CreateCompositeOptions
	DropCompositeOptions   = orm.DropCompositeOptions
	DropTableOptions       = orm.DropTableOptions
	AfterDeleteHook        = pg.AfterDeleteHook
	AfterInsertHook        = pg.AfterInsertHook
	AfterScanHook          = pg.AfterScanHook
	AfterSelectHook        = pg.AfterSelectHook
	AfterUpdateHook        = pg.AfterUpdateHook
	BeforeDeleteHook       = pg.BeforeDeleteHook
	BeforeInsertHook       = pg.BeforeInsertHook
	BeforeScanHook         = pg.BeforeScanHook
	BeforeUpdateHook       = pg.BeforeUpdateHook
)

const (
	PrimaryKeyFlag = orm.PrimaryKeyFlag
	ForeignKeyFlag = orm.ForeignKeyFlag
	NotNullFlag    = orm.NotNullFlag
	UseZeroFlag    = orm.UseZeroFlag
	UniqueFlag     = orm.UniqueFlag
	ArrayFlag      = orm.ArrayFlag
)

const (
	InvalidRelation   = orm.InvalidRelation
	HasOneRelation    = orm.HasOneRelation
	BelongsToRelation = orm.BelongsToRelation
	HasManyRelation   = orm.HasManyRelation
	Many2ManyRelation = orm.Many2ManyRelation
)

var (
	// Discard is used with Query and QueryOne to discard rows.
	Discard = pg.Discard

	// ErrMultiRows is returned by QueryOne and ExecOne when query returned multiple rows but exactly one row is expected.
	ErrMultiRows = pg.ErrMultiRows

	// ErrNoRows is returned by QueryOne and ExecOne when query returned zero rows but at least one row is expected.
	ErrNoRows = pg.ErrNoRows

	// ErrTxDone is returned by any operation that is performed on a transaction that has already been committed or rolled back.
	ErrTxDone = pg.ErrTxDone
)

var (
	// Array accepts a slice and returns a wrapper for working with PostgreSQL
	// array data type.
	//
	// For struct fields you can use array tag:
	//
	//    Emails  []string `pg:",array"`
	Array = pg.Array

	// Hstore accepts a map and returns a wrapper for working with hstore data type.
	// Supported map types are:
	//   - map[string]string
	//
	// For struct fields you can use hstore tag:
	//
	//    Attrs map[string]string `pg:",hstore"`
	Hstore = pg.Hstore

	// In accepts a slice and returns a wrapper that can be used with PostgreSQL
	// IN operator:
	//
	//    Where("id IN (?)", pg.In([]int{1, 2, 3, 4}))
	//
	// produces
	//
	//    WHERE id IN (1, 2, 3, 4)
	In = pg.In

	// InMulti accepts multiple values and returns a wrapper that can be used
	// with PostgreSQL IN operator:
	//
	//    Where("(id1, id2) IN (?)", pg.InMulti([]int{1, 2}, []int{3, 4}))
	//
	// produces
	//
	//    WHERE (id1, id2) IN ((1, 2), (3, 4))
	InMulti = pg.InMulti

	// SafeQuery replaces any placeholders found in the query.
	SafeQuery = pg.SafeQuery

	// Scan returns ColumnScanner that copies the columns in the
	// row into the values.
	Scan = pg.Scan

	// RegisterTable registers a struct as SQL table.
	// It is usually used to register intermediate table
	// in many to many relationship.
	RegisterTable = orm.RegisterTable

	// SetTableNameInflector overrides the default func that pluralizes
	// model name to get table name, e.g. my_article becomes my_articles.
	SetTableNameInflector = orm.SetTableNameInflector
)
