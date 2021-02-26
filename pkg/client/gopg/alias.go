package gopg

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type (
	Error                  = pg.Error
	CreateTableOptions     = orm.CreateTableOptions
	CreateCompositeOptions = orm.CreateCompositeOptions
	DropCompositeOptions   = orm.DropCompositeOptions
	DropTableOptions       = orm.DropTableOptions
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
