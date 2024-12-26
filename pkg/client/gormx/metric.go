package gormx

import (
	"database/sql"
	"time"

	"github.com/boxgo/box/pkg/metric"
)

type (
	Metric struct {
		ticker        *time.Ticker
		driver        string
		database      string
		statsInterval time.Duration
	}
)

const (
	labelDriver   = "driver"
	labelDatabase = "database"
	labelType     = "type"
	labelError    = "error"
)

var (
	metricConnIdle     = metric.NewGaugeVec("db_connections_idle", `The number of idle connections.`, []string{labelDriver, labelDatabase})
	metricConnInUse    = metric.NewGaugeVec("db_connections_in_use", `The number of connections currently in use.`, []string{labelDriver, labelDatabase})
	metricConnOpen     = metric.NewGaugeVec("db_connections_open", `The number of established connections both in use and idle.`, []string{labelDriver, labelDatabase})
	metricConnMaxOpen  = metric.NewGaugeVec("db_connections_max_open", `Maximum number of open connections to the database.`, []string{labelDriver, labelDatabase})
	metricWaitCount    = metric.NewGaugeVec("db_wait_count", `The total number of connections waited for.`, []string{labelDriver, labelDatabase})
	metricWaitDuration = metric.NewGaugeVec("db_wait_duration_seconds", `The total time blocked waiting for a new connection.`, []string{labelDriver, labelDatabase})
	metricSQLSeconds   = metric.NewSummaryVec("db_sql_seconds", `All queries requested seconds`, []string{labelDriver, labelDatabase, labelType, labelError}, map[float64]float64{
		0.5:  0.05,
		0.75: 0.05,
		0.9:  0.01,
		0.99: 0.001,
		1:    0.001,
	})
)

func newMetric(driver, database string, statsInterval time.Duration) *Metric {
	if statsInterval < 1 {
		statsInterval = time.Second
	}

	return &Metric{
		ticker:        time.NewTicker(statsInterval),
		driver:        driver,
		database:      database,
		statsInterval: statsInterval,
	}
}

func (m Metric) Run(db *sql.DB) error {
	go func() {
		for range m.ticker.C {
			stats := db.Stats()
			metricConnIdle.WithLabelValues(m.driver, m.database).Set(float64(stats.Idle))
			metricConnInUse.WithLabelValues(m.driver, m.database).Set(float64(stats.InUse))
			metricConnOpen.WithLabelValues(m.driver, m.database).Set(float64(stats.OpenConnections))
			metricConnMaxOpen.WithLabelValues(m.driver, m.database).Set(float64(stats.MaxOpenConnections))
			metricWaitCount.WithLabelValues(m.driver, m.database).Set(float64(stats.WaitCount))
			metricWaitDuration.WithLabelValues(m.driver, m.database).Set(stats.WaitDuration.Seconds())
		}
	}()

	return nil
}

func (m *Metric) Stop() error {
	m.ticker.Stop()
	return nil
}

func (m *Metric) registerCallback(cb *DB) error {
	if err := cb.Callback().Create().Before("gorm:before_create").Register(callbackName("before_create"), m.beforeCallback); err != nil {
		return err
	}

	if err := cb.Callback().Query().Before("gorm:before_query").Register(callbackName("before_query"), m.beforeCallback); err != nil {
		return err
	}

	if err := cb.Callback().Update().Before("gorm:before_update").Register(callbackName("before_update"), m.beforeCallback); err != nil {
		return err
	}

	if err := cb.Callback().Delete().Before("gorm:before_delete").Register(callbackName("before_delete"), m.beforeCallback); err != nil {
		return err
	}

	if err := cb.Callback().Create().After("gorm:after_create").Register(callbackName("after_create"), m.afterCallback("create")); err != nil {
		return err
	}

	if err := cb.Callback().Query().After("gorm:after_query").Register(callbackName("after_query"), m.afterCallback("query")); err != nil {
		return err
	}

	if err := cb.Callback().Update().After("gorm:after_update").Register(callbackName("after_update"), m.afterCallback("update")); err != nil {
		return err
	}

	if err := cb.Callback().Delete().After("gorm:after_delete").Register(callbackName("after_delete"), m.afterCallback("delete")); err != nil {
		return err
	}

	return nil
}

func (m *Metric) beforeCallback(db *DB) {
	db.InstanceSet("startTime", time.Now())
}

func (m *Metric) afterCallback(cmdType string) func(*DB) {
	return func(db *DB) {
		err := ""
		second := 0.0

		if db.Statement.Error != nil {
			err = db.Statement.Error.Error()
		}

		if ts, ok := db.InstanceGet("startTime"); ok {
			if startTime, ok := ts.(time.Time); ok {
				second = time.Since(startTime).Seconds()
			}
		}

		metricSQLSeconds.WithLabelValues(m.driver, m.database, cmdType, err).Observe(second)
	}
}

func callbackName(cmd string) string {
	return "gormx:" + cmd
}
