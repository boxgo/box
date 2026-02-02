package gormx

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"os"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/metric"
	"gorm.io/gorm"
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
	metricConnIdle     = metric.NewGaugeVec("db_client_connections_idle", `The number of idle connections.`, []string{labelDriver, labelDatabase})
	metricConnInUse    = metric.NewGaugeVec("db_client_connections_in_use", `The number of connections currently in use.`, []string{labelDriver, labelDatabase})
	metricConnOpen     = metric.NewGaugeVec("db_client_connections_open", `The number of established connections both in use and idle.`, []string{labelDriver, labelDatabase})
	metricConnMaxOpen  = metric.NewGaugeVec("db_client_connections_max_open", `Maximum number of open connections to the database.`, []string{labelDriver, labelDatabase})
	metricWaitCount    = metric.NewGaugeVec("db_client_connections_wait_total", `The total number of connections waited for.`, []string{labelDriver, labelDatabase})
	metricWaitDuration = metric.NewGaugeVec("db_client_connections_wait_seconds", `The total time blocked waiting for a new connection.`, []string{labelDriver, labelDatabase})
	metricSQLTotal     = metric.NewCounterVec(
		"db_client_requests_total",
		"The total number of database requests executed.",
		[]string{labelDriver, labelDatabase, labelType, "result"},
	)
	metricSQLDuration = metric.NewHistogramVec(
		"db_client_request_duration_seconds",
		"The SQL execution latencies in seconds.",
		[]string{labelDriver, labelDatabase, labelType, "result"},
		// 250us, 500us, 1ms, 2.5ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s
		[]float64{0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5},
	)
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
		second := 0.0
		result := classifyError(db.Statement.Error)

		if ts, ok := db.InstanceGet("startTime"); ok {
			if startTime, ok := ts.(time.Time); ok {
				second = time.Since(startTime).Seconds()
			}
		}

		metricSQLTotal.WithLabelValues(m.driver, m.database, cmdType, result).Inc()
		metricSQLDuration.WithLabelValues(m.driver, m.database, cmdType, result).Observe(second)
	}
}

// classifyError 将数据库错误分类为有限的几个类别，避免指标爆炸
// 同时尽可能保留有用的错误信息
func classifyError(err error) string {
	if err == nil {
		return "success"
	}

	// 检查 GORM 标准错误
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "not_found"
	}
	if errors.Is(err, gorm.ErrInvalidTransaction) {
		return "transaction_error"
	}
	if errors.Is(err, gorm.ErrMissingWhereClause) {
		return "syntax_error"
	}
	if errors.Is(err, gorm.ErrPrimaryKeyRequired) {
		return "constraint_error"
	}

	errStr := strings.ToLower(err.Error())

	// 连接相关错误
	if isConnectionError(err, errStr) {
		return "connection_error"
	}

	// 超时错误
	if isTimeoutError(err, errStr) {
		return "timeout_error"
	}

	// 约束错误（唯一键冲突、外键约束、非空约束等）
	if isConstraintError(errStr) {
		return "constraint_error"
	}

	// SQL 语法错误
	if isSyntaxError(errStr) {
		return "syntax_error"
	}

	// 事务相关错误
	if isTransactionError(errStr) {
		return "transaction_error"
	}

	// 其他错误统一归类
	return "other_error"
}

// isConnectionError 判断是否为连接相关错误
func isConnectionError(err error, errStr string) bool {
	// 检查标准库错误
	if errors.Is(err, sql.ErrConnDone) {
		return true
	}

	// 检查错误消息中的关键词
	connectionKeywords := []string{
		"connection",
		"connect",
		"connection refused",
		"connection reset",
		"connection lost",
		"connection closed",
		"no connection",
		"broken pipe",
		"network",
		"dial tcp",
		"connection timeout",
		"too many connections",
		"max connections",
		"connection pool",
		"driver: bad connection",
		"server has gone away",
		"lost connection",
	}

	for _, keyword := range connectionKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isTimeoutError 判断是否为超时错误
func isTimeoutError(err error, errStr string) bool {
	// 检查标准库超时错误
	if os.IsTimeout(err) {
		return true
	}

	// 检查 context 超时错误
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}

	// 检查 net.Error 接口的 Timeout() 方法
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	// 检查错误消息中的关键词
	timeoutKeywords := []string{
		"timeout",
		"context deadline exceeded",
		"context canceled",
		"deadline exceeded",
		"operation timed out",
		"i/o timeout",
		"read timeout",
		"write timeout",
		"query timeout",
		"statement timeout",
	}

	for _, keyword := range timeoutKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isConstraintError 判断是否为约束错误
func isConstraintError(errStr string) bool {
	constraintKeywords := []string{
		"duplicate entry",
		"unique constraint",
		"unique violation",
		"duplicate key",
		"primary key",
		"foreign key",
		"constraint violation",
		"check constraint",
		"not null",
		"cannot be null",
		"violates not-null constraint",
		"violates foreign key constraint",
		"violates unique constraint",
		"violates check constraint",
		"integrity constraint",
		"duplicate",
		"already exists",
		"1062",  // MySQL duplicate entry error code
		"23505", // PostgreSQL unique violation error code
		"23503", // PostgreSQL foreign key violation error code
		"23502", // PostgreSQL not null violation error code
	}

	for _, keyword := range constraintKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isSyntaxError 判断是否为 SQL 语法错误
func isSyntaxError(errStr string) bool {
	syntaxKeywords := []string{
		"syntax error",
		"sql syntax",
		"parse error",
		"invalid syntax",
		"unexpected token",
		"unexpected end",
		"missing",
		"unknown column",
		"unknown table",
		"table doesn't exist",
		"column doesn't exist",
		"1064",  // MySQL syntax error code
		"42601", // PostgreSQL syntax error code
	}

	for _, keyword := range syntaxKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isTransactionError 判断是否为事务相关错误
func isTransactionError(errStr string) bool {
	transactionKeywords := []string{
		"transaction",
		"deadlock",
		"lock wait timeout",
		"lock wait",
		"could not serialize",
		"serialization failure",
		"transaction rollback",
		"transaction commit",
		"in failed sql transaction",
		"current transaction is aborted",
	}

	for _, keyword := range transactionKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

func callbackName(cmd string) string {
	return "gormx:" + cmd
}
