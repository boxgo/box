package gormx

import (
	"context"

	"github.com/boxgo/box/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type (
	Gorm struct {
		cfg    *Config
		db     *gorm.DB
		metric *Metric
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

	var namer schema.Namer
	if c.NamingStrategy != nil {
		namer = c.NamingStrategy
	} else {
		namer = schema.NamingStrategy{
			TablePrefix:   c.NamingStrategyTablePrefix,
			SingularTable: c.NamingStrategySingularTable,
			NameReplacer:  c.NamingStrategyNameReplacer,
			NoLowerCase:   c.NamingStrategyNoLowerCase,
		}
	}

	db, err := gorm.Open(dial, &gorm.Config{
		SkipDefaultTransaction:                   c.SkipDefaultTransaction,
		DryRun:                                   c.DryRun,
		PrepareStmt:                              c.PrepareStmt,
		DisableNestedTransaction:                 c.DisableNestedTransaction,
		AllowGlobalUpdate:                        c.AllowGlobalUpdate,
		DisableForeignKeyConstraintWhenMigrating: c.DisableForeignKeyConstraintWhenMigrating,
		QueryFields:                              c.QueryFields,
		CreateBatchSize:                          c.CreateBatchSize,
		NamingStrategy:                           namer,
		NowFunc:                                  c.NowFunc,
		ConnPool:                                 c.ConnPool,
		ClauseBuilders:                           c.ClauseBuilders,
		Plugins:                                  c.Plugins,
		Logger:                                   &Logger{},
	})
	if err != nil {
		logger.Panicf("Gorm open error %s %s: %s", c.Driver, c.DSN, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panicf("Gorm get db error %s %s: %s", c.Driver, c.DSN, err)
	}
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(c.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(c.MaxLifeTime)

	metric := newMetric(c.Driver, c.Key(), c.MetricInterval)
	if err := metric.registerCallback(db); err != nil {
		logger.Panicf("metric register callback error: %s", err)
	}

	return &Gorm{
		cfg:    c,
		db:     db,
		metric: metric,
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

	if err := orm.metric.Run(db); err != nil {
		return err
	}

	return db.Ping()
}

func (orm *Gorm) Shutdown(ctx context.Context) error {
	db, err := orm.db.DB()
	if err != nil {
		return err
	}

	if err = orm.metric.Stop(); err != nil {
		return err
	}

	return db.Close()
}

func (orm Gorm) DB() *DB {
	return orm.db
}
