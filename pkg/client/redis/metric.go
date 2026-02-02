package redis

import (
	"context"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/metric"
	"github.com/redis/go-redis/v9"
)

type (
	Metric struct {
		cfg  *Config
		addr string
	}

	startKey struct{}
)

func newMetric(cfg *Config) *Metric {
	return &Metric{
		cfg:  cfg,
		addr: strings.Join(cfg.Address, ","),
	}
}

var (
	cmdTotal = metric.NewCounterVec(
		"redis_client_requests_total",
		"The total number of Redis commands executed.",
		[]string{"pipe", "cmd", "result"},
	)
	cmdDuration = metric.NewHistogramVec(
		"redis_client_request_duration_seconds",
		"The Redis command latencies in seconds.",
		[]string{"pipe", "cmd", "result"},
		// 100us, 250us, 500us, 1ms, 2.5ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms
		[]float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5},
	)
)

func (m *Metric) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (m *Metric) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmd)

		m.report(ctx, false, time.Since(start), cmd)

		return err
	}
}

func (m *Metric) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmds)

		m.report(ctx, false, time.Since(start), cmds...)

		return err
	}
}

func (m *Metric) report(ctx context.Context, pipe bool, elapsed time.Duration, cmds ...redis.Cmder) {
	cmdStr := ""
	result := "success"
	pipeStr := strconv.FormatBool(pipe)

	if pipe {
		cmdStr = "pipeline"
	} else if len(cmds) > 0 {
		cmdStr = cmds[0].Name()
	}

	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil && err != redis.Nil {
			result = classifyRedisError(err)
			break
		}
	}

	values := []string{
		pipeStr,
		cmdStr,
		result,
	}

	cmdDuration.WithLabelValues(values...).Observe(elapsed.Seconds())
	cmdTotal.WithLabelValues(values...).Inc()
}

// classifyRedisError 将 Redis 错误分类为有限的几个类别，避免指标爆炸
// 同时尽可能保留有用的错误信息
func classifyRedisError(err error) string {
	if err == nil {
		return "success"
	}

	// 检查 redis.Nil（键不存在，这是正常情况，不应该算作错误）
	if err == redis.Nil {
		return "success"
	}

	// 检查标准库超时错误
	if os.IsTimeout(err) {
		return "timeout_error"
	}

	// 检查 context 超时错误
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
		return "timeout_error"
	}

	// 检查 net.Error 接口的 Timeout() 方法
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "timeout_error"
	}

	// 检查连接相关错误
	if isRedisConnectionError(err) {
		return "connection_error"
	}

	// 检查事务错误
	if err == redis.TxFailedErr || errors.Is(err, redis.TxFailedErr) {
		return "transaction_error"
	}

	errStr := strings.ToLower(err.Error())

	// 检查 Redis 命令错误
	if isRedisCommandError(errStr) {
		return "command_error"
	}

	// 检查事务相关错误
	if isRedisTransactionError(errStr) {
		return "transaction_error"
	}

	// 检查权限错误
	if isRedisAuthError(errStr) {
		return "auth_error"
	}

	// 检查内存不足错误
	if isRedisOOMError(errStr) {
		return "oom_error"
	}

	// 检查集群相关错误
	if isRedisClusterError(errStr) {
		return "cluster_error"
	}

	// 其他错误统一归类
	return "other_error"
}

// isRedisConnectionError 判断是否为连接相关错误
func isRedisConnectionError(err error) bool {
	// 检查标准库错误
	if errors.Is(err, redis.ErrClosed) {
		return true
	}

	errStr := strings.ToLower(err.Error())

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
		"i/o error",
		"use of closed network connection",
		"connection reset by peer",
		"no such host",
		"no route to host",
		"refused",
		"closed",
		"EOF",
	}

	for _, keyword := range connectionKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isRedisCommandError 判断是否为 Redis 命令错误
func isRedisCommandError(errStr string) bool {
	commandKeywords := []string{
		"wrongtype",
		"wrong type",
		"wrong number of arguments",
		"unknown command",
		"command not allowed",
		"invalid argument",
		"invalid command",
		"syntax error",
		"parse error",
		"protocol error",
		"ERR", // Redis 错误前缀
		"WRONGTYPE",
		"NOSCRIPT", // Lua 脚本不存在
		"BUSYKEY",  // 键正在被其他操作使用
	}

	for _, keyword := range commandKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isRedisTransactionError 判断是否为事务相关错误
func isRedisTransactionError(errStr string) bool {
	transactionKeywords := []string{
		"transaction",
		"EXECABORT",
		"transaction failed",
		"watch",
		"CAS", // Compare and Swap
	}

	for _, keyword := range transactionKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isRedisAuthError 判断是否为权限/认证错误
func isRedisAuthError(errStr string) bool {
	authKeywords := []string{
		"noauth",
		"authentication required",
		"invalid password",
		"auth",
		"permission denied",
		"ACL",
	}

	for _, keyword := range authKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isRedisOOMError 判断是否为内存不足错误
func isRedisOOMError(errStr string) bool {
	oomKeywords := []string{
		"oom",
		"out of memory",
		"command not allowed when used memory",
		"maxmemory",
	}

	for _, keyword := range oomKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// isRedisClusterError 判断是否为集群相关错误
func isRedisClusterError(errStr string) bool {
	clusterKeywords := []string{
		"cluster",
		"MOVED",
		"ASK",
		"CLUSTERDOWN",
		"TRYAGAIN",
		"crossslot",
		"slot",
		"migrating",
		"importing",
	}

	for _, keyword := range clusterKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}
