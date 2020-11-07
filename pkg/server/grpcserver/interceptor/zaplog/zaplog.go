package zaplog

import (
	"context"

	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/system"
	"github.com/boxgo/box/pkg/util/strutil"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		newCtx := wrapCtx(ctx, md, info.FullMethod)

		logger.TraceRaw(newCtx).Info(">>>", []zap.Field{zap.Any("req", req)}...)

		resp, err = handler(newCtx, req)
		if err != nil {
			logger.TraceRaw(newCtx).Info("xxx", []zap.Field{zap.Any("err", err)}...)
		}

		logger.TraceRaw(newCtx).Info("<<<", []zap.Field{zap.Any("resp", resp)}...)

		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, _ := metadata.FromIncomingContext(ctx)
		newCtx := wrapCtx(ctx, md, info.FullMethod)

		logger.TraceRaw(newCtx).Info(">>>")

		err := handler(newCtx, ss)
		if err != nil {
			logger.TraceRaw(newCtx).Info("xxx", []zap.Field{zap.Any("err", err)}...)
		}

		logger.TraceRaw(newCtx).Info("<<<", []zap.Field{zap.Any("resp", "")}...)

		return err
	}
}

func wrapCtx(ctx context.Context, md metadata.MD, biz string) context.Context {
	uidKey := system.TraceUID()
	reqIdKey := system.TraceReqID()
	bizIdKey := system.TraceBizID()
	spanKey := system.TraceSpanID()

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
