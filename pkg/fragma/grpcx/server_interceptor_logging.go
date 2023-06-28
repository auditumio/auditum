package grpcx

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggingUnaryServerInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return grpc_zap.UnaryServerInterceptor(log.Named("grpc"))
}
