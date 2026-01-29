package wukong

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"
)

func TestClassifyHTTPError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "success",
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: "timeout_error",
		},
		{
			name:     "os deadline exceeded",
			err:      os.ErrDeadlineExceeded,
			expected: "timeout_error",
		},
		{
			name:     "os timeout",
			err:      &os.SyscallError{Err: os.ErrDeadlineExceeded},
			expected: "timeout_error",
		},
		{
			name: "net timeout error",
			err: &net.OpError{
				Op:  "dial",
				Err: &os.SyscallError{Err: os.ErrDeadlineExceeded},
			},
			expected: "timeout_error",
		},
		{
			name:     "dns error",
			err:      errors.New("no such host: example.com"),
			expected: "dns_error",
		},
		{
			name:     "dns lookup error",
			err:      errors.New("lookup example.com: no such host"),
			expected: "dns_error",
		},
		{
			name:     "tls error",
			err:      errors.New("tls: handshake failure"),
			expected: "tls_error",
		},
		{
			name:     "tls certificate error",
			err:      errors.New("x509: certificate verify failed"),
			expected: "tls_error",
		},
		{
			name:     "connection refused",
			err:      errors.New("connection refused"),
			expected: "connection_error",
		},
		{
			name:     "connection reset",
			err:      errors.New("connection reset by peer"),
			expected: "connection_error",
		},
		{
			name:     "network error",
			err:      errors.New("dial tcp 127.0.0.1:8080: connect: connection refused"),
			expected: "connection_error",
		},
		{
			name:     "EOF error",
			err:      errors.New("EOF"),
			expected: "connection_error",
		},
		{
			name:     "net error without timeout",
			err:      &net.OpError{Op: "dial", Err: errors.New("connection refused")},
			expected: "connection_error",
		},
		{
			name:     "other error",
			err:      errors.New("some unknown error"),
			expected: "other_error",
		},
		{
			name:     "http error message",
			err:      errors.New("bad request"),
			expected: "other_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyHTTPError(tt.err)
			if result != tt.expected {
				t.Errorf("classifyHTTPError(%v) = %q, want %q", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsDNSError(t *testing.T) {
	tests := []struct {
		name     string
		errStr   string
		expected bool
	}{
		{"no such host", "no such host: example.com", true},
		{"lookup error", "lookup example.com: no such host", true},
		{"unknown host", "unknown host", true},
		{"host not found", "host not found", true},
		{"not dns error", "connection refused", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDNSError(tt.errStr)
			if result != tt.expected {
				t.Errorf("isDNSError(%q) = %v, want %v", tt.errStr, result, tt.expected)
			}
		})
	}
}

func TestIsTLSError(t *testing.T) {
	tests := []struct {
		name     string
		errStr   string
		expected bool
	}{
		{"tls handshake", "tls: handshake failure", true},
		{"certificate error", "x509: certificate verify failed", true},
		{"ssl error", "ssl handshake failure", true},
		{"not tls error", "connection refused", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTLSError(tt.errStr)
			if result != tt.expected {
				t.Errorf("isTLSError(%q) = %v, want %v", tt.errStr, result, tt.expected)
			}
		})
	}
}

func TestIsHTTPConnectionError(t *testing.T) {
	tests := []struct {
		name     string
		errStr   string
		expected bool
	}{
		{"connection refused", "connection refused", true},
		{"connection reset", "connection reset by peer", true},
		{"dial tcp", "dial tcp 127.0.0.1:8080: connect: connection refused", true},
		{"EOF", "EOF", true},
		{"network unreachable", "network is unreachable", true},
		{"not connection error", "dns lookup failed", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHTTPConnectionError(tt.errStr)
			if result != tt.expected {
				t.Errorf("isHTTPConnectionError(%q) = %v, want %v", tt.errStr, result, tt.expected)
			}
		})
	}
}
