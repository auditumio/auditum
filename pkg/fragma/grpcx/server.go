// Copyright 2023 Igor Zibarev
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpcx

import (
	"runtime/debug"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
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
	panicRecoveryInterceptor := recovery.UnaryServerInterceptor(
		recovery.WithRecoveryHandler(panicRecoveryHandler),
	)

	// NOTE: two interceptors is not a mistake. The last one should catch RPC
	// panics, so interceptors can build based on it, e.g. finish call logging.
	// The first one should catch all other panics, i.e. it is a guard for
	// other interceptor panics.

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			panicRecoveryInterceptor,
			LoggingUnaryServerInterceptor(log),
			otelgrpc.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			panicRecoveryInterceptor,
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

func newPanicRecoveryHandler(log *zap.Logger) func(p any) error {
	return func(p any) error {
		log.Error("Recover from panic",
			zap.Any("panic", p),
			zap.String("stack", string(debug.Stack())),
		)
		return status.Error(codes.Internal, "")
	}
}
