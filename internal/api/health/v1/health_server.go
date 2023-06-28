package healthv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthServer implements gRPC health server.
type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
}

func NewHealthServer() *HealthServer {
	return &HealthServer{}
}

func (s *HealthServer) Check(
	context.Context,
	*grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthServer) Register(srv *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(srv, s)
}
