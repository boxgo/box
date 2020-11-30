package zaplog

import (
	"context"

	"github.com/boxgo/box/pkg/grpc/wrapper"
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/trace"
	"github.com/boxgo/box/pkg/util/strutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		newCtx := wrapCtx(ctx, md, info.FullMethod)

		logger.Trace(newCtx).Infow(">>>", "req", req)

		resp, err = handler(newCtx, req)
		if err != nil {
			logger.Trace(newCtx).Infow("xxx", "err", err)
		}

		logger.Trace(newCtx).Infow("<<<", "resp", resp)

		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, _ := metadata.FromIncomingContext(ctx)

		wrapped := wrapper.WrapServerStream(ss)
		wrapped.WrappedContext = wrapCtx(ctx, md, info.FullMethod)

		logger.Trace(wrapped.WrappedContext).Info(">>>")

		err := handler(srv, wrapped)
		if err != nil {
			logger.Trace(wrapped.WrappedContext).Infow("xxx", "err", err)
		}

		logger.Trace(wrapped.WrappedContext).Info("<<<")

		return err
	}
}

func wrapCtx(ctx context.Context, md metadata.MD, biz string) context.Context {
	uidKey := trace.ID()
	reqIdKey := trace.ReqID()
	bizIdKey := trace.BizID()
	spanKey := trace.SpanID()

	uidVal := strutil.First(md.Get(uidKey))
	reqIdVal := strutil.First(md.Get(reqIdKey))
	spanVal := strutil.First(md.Get(spanKey))
	bizIdVal := biz

	newCtx := context.WithValue(ctx, uidKey, uidVal)
	newCtx = context.WithValue(newCtx, reqIdKey, reqIdVal)
	newCtx = context.WithValue(newCtx, bizIdKey, bizIdVal)
	newCtx = context.WithValue(newCtx, spanKey, spanVal)

	newMd := md.Copy()
	newMd.Set(uidKey, uidVal)
	newMd.Set(uidKey, uidVal)
	newMd.Set(reqIdKey, reqIdVal)
	newMd.Set(bizIdKey, bizIdVal)
	newMd.Set(spanKey, spanVal)

	return metadata.NewIncomingContext(newCtx, newMd)
}
