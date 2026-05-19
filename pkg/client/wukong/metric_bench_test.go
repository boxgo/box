package wukong

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"testing"
)

// 模拟各种错误类型用于性能测试
var (
	testErrors = []struct {
		name string
		err  error
	}{
		{"nil", nil},
		{"context_deadline", context.DeadlineExceeded},
		{"os_timeout", &os.SyscallError{Err: os.ErrDeadlineExceeded}},
		{"net_timeout", &net.OpError{Err: &os.SyscallError{Err: os.ErrDeadlineExceeded}}},
		{"dns_error", errors.New("no such host: example.com")},
		{"tls_error", errors.New("tls: handshake failure")},
		{"connection_error", errors.New("connection refused")},
		{"other_error", errors.New("some unknown error")},
	}
)

// BenchmarkClassifyHTTPError 测试错误分类函数的性能
func BenchmarkClassifyHTTPError(b *testing.B) {
	b.Run("success_case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			classifyHTTPError(nil)
		}
	})

	b.Run("timeout_error", func(b *testing.B) {
		err := context.DeadlineExceeded
		for i := 0; i < b.N; i++ {
			classifyHTTPError(err)
		}
	})

	b.Run("os_timeout", func(b *testing.B) {
		err := &os.SyscallError{Err: os.ErrDeadlineExceeded}
		for i := 0; i < b.N; i++ {
			classifyHTTPError(err)
		}
	})

	b.Run("net_timeout", func(b *testing.B) {
		err := &net.OpError{Err: &os.SyscallError{Err: os.ErrDeadlineExceeded}}
		// 需要设置 Timeout() 方法返回 true
		for i := 0; i < b.N; i++ {
			classifyHTTPError(err)
		}
	})

	b.Run("connection_error", func(b *testing.B) {
		err := errors.New("connection refused")
		for i := 0; i < b.N; i++ {
			classifyHTTPError(err)
		}
	})

	b.Run("dns_error", func(b *testing.B) {
		err := errors.New("no such host: example.com")
		for i := 0; i < b.N; i++ {
			classifyHTTPError(err)
		}
	})

	b.Run("tls_error", func(b *testing.B) {
		err := errors.New("tls: handshake failure")
		for i := 0; i < b.N; i++ {
			classifyHTTPError(err)
		}
	})

	b.Run("other_error", func(b *testing.B) {
		err := errors.New("some unknown error")
		for i := 0; i < b.N; i++ {
			classifyHTTPError(err)
		}
	})

	b.Run("mixed_errors", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			testCase := testErrors[i%len(testErrors)]
			classifyHTTPError(testCase.err)
		}
	})
}

// BenchmarkClassifyHTTPError_Old 测试旧的简单错误处理（作为对比）
func BenchmarkClassifyHTTPError_Old(b *testing.B) {
	b.Run("simple_error_check", func(b *testing.B) {
		err := errors.New("some error")
		for i := 0; i < b.N; i++ {
			if err != nil {
				_ = "error"
			}
		}
	})
}

// BenchmarkStringOperations 测试字符串操作的开销
func BenchmarkStringOperations(b *testing.B) {
	err := errors.New("connection refused: dial tcp 127.0.0.1:8080")

	b.Run("error_string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = err.Error()
		}
	})

	b.Run("to_lower", func(b *testing.B) {
		errStr := err.Error()
		for i := 0; i < b.N; i++ {
			_ = strings.ToLower(errStr)
		}
	})

	b.Run("contains_check", func(b *testing.B) {
		errStr := err.Error()
		keywords := []string{"connection", "refused", "dial", "tcp"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, keyword := range keywords {
				if strings.Contains(errStr, keyword) {
					break
				}
			}
		}
	})
}
