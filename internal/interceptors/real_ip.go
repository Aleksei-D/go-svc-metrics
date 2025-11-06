package interceptors

import (
	"context"
	"go-svc-metrics/internal/logger"
	"net"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewRealIPInterceptor(network *net.IPNet) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if isTrustedNetwork(md, network) {
				return handler(ctx, req)
			}
		}
		return nil, status.Error(codes.Unauthenticated, "invalid real IP")
	}
}

func isTrustedNetwork(md metadata.MD, network *net.IPNet) bool {
	ipData := md.Get("X-Real-IP")
	if len(ipData) <= 0 {
		return false
	}

	parts := strings.Split(ipData[0], ",")
	if len(parts) <= 0 {
		return false
	}

	ipSTR := strings.TrimSpace(parts[0])
	ip := net.ParseIP(ipSTR)
	if ip == nil {
		return false
	}

	if !network.Contains(ip) {
		return false
	}

	return true
}

func NewRealIPClientInterceptor(realIP string) func(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md := metadata.New(map[string]string{"X-Real-IP": realIP})
		newCtx := metadata.NewOutgoingContext(ctx, md)

		err := invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			logger.Log.Warn("REQUEST", zap.String("URI", method), zap.Error(err))
		}
		return err
	}
}
