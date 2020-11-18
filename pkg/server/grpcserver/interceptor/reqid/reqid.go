package reqid

import (
	"context"

	"github.com/boxgo/box/pkg/grpc/wrapper"
	"github.com/boxgo/box/pkg/trace"
	"github.com/boxgo/box/pkg/util/strutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryRequestId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		newCtx := wrapCtx(ctx, md)

		return handler(newCtx, req)
	}
}

func StreamRequestId() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(srv, ss)
		}

		wrapped := wrapper.WrapServerStream(ss)
		wrapped.WrappedContext = wrapCtx(ctx, md)

		return handler(srv, wrapped)
	}
}

func wrapCtx(ctx context.Context, md metadata.MD) context.Context {
	key := trace.ReqID()
	val, exist := getReqId(md)

	if exist {
		return ctx
	}

	newMd := md.Copy()
	newMd.Set(key, val)

	return metadata.NewIncomingContext(ctx, newMd)
}

func getReqId(md metadata.MD) (string, bool) {
	if id := strutil.First(md.Get(trace.ReqID())); id != "" {
		return id, true
	}

	return strutil.ShortID(), false
}
