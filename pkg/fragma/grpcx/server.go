package grpcx

import (
	"runtime/debug"
	"time"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func NewServer(log *zap.Logger) *grpc.Server {
	panicRecoveryHandler := newPanicRecoveryHandler(log.Named("grpc_panic_handler"))

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(panicRecoveryHandler),
			),
			otelgrpc.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			LoggingUnaryServerInterceptor(log),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      10 * time.Minute,
			MaxConnectionAgeGrace: 1 * time.Minute,
			Time:                  2 * time.Minute,
			Timeout:               1 * time.Minute,
		}),
	)

	reflection.Register(server)

	return server
}

func InitPrometheusMetrics(server *grpc.Server) {
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(server)
}

func newPanicRecoveryHandler(log *zap.Logger) func(p interface{}) error {
	return func(p interface{}) error {
		log.Error("Recover from panic",
			zap.Any("panic", p),
			zap.String("stack", string(debug.Stack())),
		)
		return status.Error(codes.Internal, "")
	}
}
