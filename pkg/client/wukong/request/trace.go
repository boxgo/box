package request

import (
	"log"
	"net/http/httptrace"
)

var (
	trace = &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.Println("GetConn", hostPort)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Println("GotConn", connInfo)
		},
		GotFirstResponseByte: func() {
			log.Println("GotFirstResponseByte")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.Println("DNSStart", info)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			log.Println("DNSDone", dnsInfo)
		},
		ConnectStart: func(network, addr string) {
			log.Println("ConnectStart", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			log.Println("ConnectDone", network, addr, err)
		},
	}
)
