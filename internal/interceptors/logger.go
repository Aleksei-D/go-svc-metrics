package interceptors

import (
	"context"
	"go-svc-metrics/internal/logger"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)

	duration := time.Since(start)
	logger.Log.Info("REQUEST", zap.String("URI", info.FullMethod), zap.Duration("DURATION", duration))
	switch err {
	case nil:
		size := proto.Size(resp.(protoreflect.ProtoMessage))
		logger.Log.Info("RESPONSE", zap.String("URI", info.FullMethod), zap.Int("SIZE", size))
	default:
		logger.Log.Warn("RESPONSE", zap.String("URI", info.FullMethod), zap.Error(err), zap.Duration("DURATION", duration))
	}

	return resp, err
}
