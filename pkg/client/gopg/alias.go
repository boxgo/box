package gopg

import (
	"github.com/go-pg/pg/v10/orm"
)

type (
	CreateTableOptions     = orm.CreateTableOptions
	CreateCompositeOptions = orm.CreateCompositeOptions
	DropCompositeOptions   = orm.DropCompositeOptions
	DropTableOptions       = orm.DropTableOptions
)
