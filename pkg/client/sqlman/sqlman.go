package sqlman

import (
	"database/sql"

	"github.com/boxgo/box/pkg/logger"
)

type (
	SQLMan struct {
		cfg *Config
		db  *sql.DB
	}
)

func newSQLMan(c *Config) *SQLMan {
	db, err := sql.Open(c.Driver, c.DSN)
	if err != nil {
		logger.Panicf("SQLMan open %s %s error: %s", c.Driver, err)
	}

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxIdleTime(c.MaxIdleTime)
	db.SetConnMaxLifetime(c.MaxLifeTime)

	return &SQLMan{
		cfg: c,
		db:  db,
	}
}

func (s SQLMan) DB() *sql.DB {
	return s.db
}
