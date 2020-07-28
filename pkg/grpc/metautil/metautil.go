package metautil

import (
	"context"
	"encoding/base64"
	"strings"

	"google.golang.org/grpc/metadata"
)

type NewMD metadata.MD

const (
	binHdrSuffix = "-bin"
)

func ExtractIncoming(ctx context.Context) NewMD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return NewMD(metadata.Pairs())
	}
	return NewMD(md)
}

func ExtractOutgoing(ctx context.Context) NewMD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return NewMD(metadata.Pairs())
	}
	return NewMD(md)
}

func (m NewMD) ToOutgoing(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD(m))
}

func (m NewMD) ToIncoming(ctx context.Context) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.MD(m))
}

func (m NewMD) Get(key string) string {
	k, _ := encodeKeyValue(key, "")
	vv, ok := m[k]
	if !ok {
		return ""
	}
	return vv[0]
}

func (m NewMD) Del(key string) NewMD {
	k, _ := encodeKeyValue(key, "")
	delete(m, k)
	return m
}

func (m NewMD) Set(key string, value string) NewMD {
	k, v := encodeKeyValue(key, value)
	m[k] = []string{v}
	return m
}

func (m NewMD) Add(key string, value string) NewMD {
	k, v := encodeKeyValue(key, value)
	m[k] = append(m[k], v)
	return m
}

func encodeKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, binHdrSuffix) {
		val := base64.StdEncoding.EncodeToString([]byte(v))
		v = string(val)
	}
	return k, v
}
