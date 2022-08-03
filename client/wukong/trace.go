package wukong

import (
	"crypto/tls"
	"net/http/httptrace"
	"time"
)

type (
	TraceInfo struct {
		HostPort               string
		GotConnInfo            httptrace.GotConnInfo
		DNSStartInfo           httptrace.DNSStartInfo
		DNSDoneInfo            httptrace.DNSDoneInfo
		ConnectNetwork         string
		ConnectAddr            string
		ConnectError           error
		TLSConnectionState     tls.ConnectionState
		TLSError               error
		ElapsedTime            time.Duration
		GetConnAt              time.Time
		GotConnAt              time.Time
		GetConnElapsed         time.Duration
		DNSStartAt             time.Time
		DNSDoneAt              time.Time
		DNSLookupElapsed       time.Duration
		ConnectStartAt         time.Time
		ConnectDoneAt          time.Time
		ConnectElapsed         time.Duration
		TLSHandshakeStartAt    time.Time
		TLSHandshakeDoneAt     time.Time
		TLSHandshakeElapsed    time.Duration
		GotFirstResponseByteAt time.Time
	}
)

func traceGenerator(request *Request) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			request.TraceInfo.GetConnAt = time.Now()
			request.TraceInfo.HostPort = hostPort
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			request.TraceInfo.GotConnAt = time.Now()
			request.TraceInfo.GetConnElapsed = request.TraceInfo.GotConnAt.Sub(request.TraceInfo.GetConnAt)
			request.TraceInfo.GotConnInfo = connInfo
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			request.TraceInfo.DNSStartAt = time.Now()
			request.TraceInfo.DNSStartInfo = info
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			request.TraceInfo.DNSDoneAt = time.Now()
			request.TraceInfo.DNSLookupElapsed = request.TraceInfo.DNSDoneAt.Sub(request.TraceInfo.DNSStartAt)
			request.TraceInfo.DNSDoneInfo = dnsInfo
		},
		ConnectStart: func(network, addr string) {
			request.TraceInfo.ConnectStartAt = time.Now()
			request.TraceInfo.ConnectAddr = addr
			request.TraceInfo.ConnectNetwork = network
		},
		ConnectDone: func(network, addr string, err error) {
			request.TraceInfo.ConnectDoneAt = time.Now()
			request.TraceInfo.ConnectElapsed = request.TraceInfo.ConnectDoneAt.Sub(request.TraceInfo.ConnectStartAt)
			request.TraceInfo.ConnectError = err
		},
		TLSHandshakeStart: func() {
			request.TraceInfo.TLSHandshakeStartAt = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			request.TraceInfo.TLSHandshakeDoneAt = time.Now()
			request.TraceInfo.TLSHandshakeElapsed = request.TraceInfo.TLSHandshakeDoneAt.Sub(request.TraceInfo.TLSHandshakeStartAt)
			request.TraceInfo.TLSConnectionState = state
		},
		GotFirstResponseByte: func() {
			request.TraceInfo.GotFirstResponseByteAt = time.Now()
		},
	}
}
