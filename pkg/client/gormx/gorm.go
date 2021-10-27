package gormx

import (
	"context"
	"database/sql"

	"github.com/boxgo/box/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	Gorm struct {
		cfg   *Config
		db    *gorm.DB
		rawDB *sql.DB
	}
)

func newGorm(c *Config) *Gorm {
	var dial gorm.Dialector

	if c.dial != nil {
		dial = c.dial
	} else {
		switch c.Driver {
		case "mysql":
			dial = mysql.Open(c.DSN)
		case "postgres":
			dial = postgres.Open(c.DSN)
		default:
			logger.Panicf("Gorm not support %s", c.DSN)
		}
	}

	db, err := gorm.Open(dial, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 &Logger{},
	})
	if err != nil {
		logger.Panicf("Gorm open error %s %s: %s", c.Driver, c.DSN, err)
	}

	rawDb, err := db.DB()
	if err != nil {
		logger.Panicf("Gorm get db error %s %s: %s", c.Driver, c.DSN, err)
	}

	rawDb.SetMaxOpenConns(c.MaxOpenConns)
	rawDb.SetMaxIdleConns(c.MaxIdleConns)
	rawDb.SetConnMaxIdleTime(c.MaxIdleTime)
	rawDb.SetConnMaxLifetime(c.MaxLifeTime)

	return &Gorm{
		cfg:   c,
		db:    db,
		rawDB: rawDb,
	}
}

func (orm Gorm) Name() string {
	return "gorm"
}

func (orm *Gorm) Serve(ctx context.Context) error {
	db, err := orm.db.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}

func (orm *Gorm) Shutdown(ctx context.Context) error {
	db, err := orm.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (orm Gorm) DB() *DB {
	return orm.db
}
