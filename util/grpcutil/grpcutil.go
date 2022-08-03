package grpcutil

import (
	"context"
	"net"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// GetRealAddr get x-real-ip
// nginx config eg: "rpc_set_header X-Real-IP $remote_addr;"
func GetRealAddr(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	rips := md.Get("x-real-ip")
	if len(rips) == 0 {
		return ""
	}

	return rips[0]
}

// GetPeerAddr get peer addr
func GetPeerAddr(ctx context.Context) string {
	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
		} else {
			addr = pr.Addr.String()
		}
	}
	return addr
}

// GetRealPeerAddr returns real ip when exist, otherwise returns peer ip.
func GetRealPeerAddr(ctx context.Context) string {
	if ip := GetRealAddr(ctx); ip != "" {
		return ip
	}

	return GetPeerAddr(ctx)
}
