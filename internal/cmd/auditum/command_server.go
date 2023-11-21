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

package auditum

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/uptrace/bun"
	"go.uber.org/zap"

	auditumv1alpha1 "github.com/auditumio/auditum/internal/api/auditumio/auditum/v1alpha1"
	healthv1 "github.com/auditumio/auditum/internal/api/health/v1"
	"github.com/auditumio/auditum/internal/aud"
	"github.com/auditumio/auditum/internal/grpcgateway"
	"github.com/auditumio/auditum/internal/sql"
	"github.com/auditumio/auditum/internal/sql/postgres"
	"github.com/auditumio/auditum/internal/sql/sqlite"
	"github.com/auditumio/auditum/pkg/fragma/bunx"
	"github.com/auditumio/auditum/pkg/fragma/grpcx"
	"github.com/auditumio/auditum/pkg/fragma/httpx"
	"github.com/auditumio/auditum/pkg/fragma/otelx"
	"github.com/auditumio/auditum/pkg/fragma/uds"
)

func executeServer(conf *Configuration, log *zap.Logger) int {
	// --- Startup phase ---

	ctx := context.Background()

	slog := log.Sugar()

	slog.Infof("%s %s is starting...", appName, commandNameServer)

	tracingProvider, err := initTracing(conf.Tracing, log)
	if err != nil {
		log.Error("Failed to initialize tracing", zap.Error(err))
		return exitCodeStartFailure
	}

	settings := aud.Settings{
		Records: conf.Settings.Records,
	}

	var db *bun.DB
	switch conf.Store.Type {
	case storeTypeSQLite:
		var err error
		db, err = sqlite.NewDatabase(
			ctx,
			conf.Store.SQLite.DatabasePath,
			log,
			bunx.LogQueriesFlagFromBool(conf.Store.SQLite.LogQueries),
		)
		if err != nil {
			log.Error("Failed to connect to database", zap.Error(err))
			return exitCodeStartFailure
		}

		if conf.Store.SQLite.DatabasePath == sqlite.FilepathMemory {
			log.Warn("Using in-memory SQLite database. All data will be lost on shutdown.")

			if err := sqlite.RunMigrations(
				db,
				conf.Store.SQLite.DatabasePath,
				conf.Store.SQLite.MigrationsPath,
				log,
			); err != nil {
				log.Error("Failed to run migrations", zap.Error(err))
				return exitCodeStartFailure
			}
		}
	case storeTypePostgres:
		var err error
		db, err = postgres.NewDatabase(
			ctx,
			conf.Store.Postgres.Host,
			conf.Store.Postgres.Port,
			conf.Store.Postgres.Database,
			conf.Store.Postgres.Username,
			conf.Store.Postgres.Password,
			conf.Store.Postgres.SSLMode,
			log,
			bunx.LogQueriesFlagFromBool(conf.Store.Postgres.LogQueries),
		)
		if err != nil {
			log.Error("Failed to connect to database", zap.Error(err))
			return exitCodeStartFailure
		}
	default:
		log.Panic("Unreachable code: invalid store type", zap.String("store_type", conf.Store.Type))
		return exitCodeStartFailure
	}

	store := sql.NewStore(db)

	unixSocketAvailable := true
	if err := uds.IsAvailable(); err != nil {
		log.Warn(
			"Unix socket is not available. "+
				"gRPC gateway will connect to gRPC server via TCP.",
			zap.Error(err),
		)
		unixSocketAvailable = false
	}

	unixSocket := ""
	if unixSocketAvailable {
		unixSocket, err = uds.NewSocket()
		if err != nil {
			log.Error("Failed to create UNIX socket", zap.Error(err))
			return exitCodeStartFailure
		}
	}

	grpcServer := grpcx.NewServer(log)

	healthServer := healthv1.NewHealthServer()
	healthServer.Register(grpcServer)

	projectServiceServer := auditumv1alpha1.NewProjectServiceServer(
		store,
		log,
	)
	projectServiceServer.RegisterServer(grpcServer)

	recordServiceServer := auditumv1alpha1.NewRecordServiceServer(
		store,
		log,
		settings,
	)
	recordServiceServer.RegisterServer(grpcServer)

	// NOTE: must be called after all services are registered.
	grpcx.InitPrometheusMetrics(grpcServer)

	grpcServerAddr := ":" + conf.GRPC.Port

	var grpcServerControllerOpts []grpcx.ServerControllerOption
	if unixSocketAvailable {
		grpcServerControllerOpts = append(
			grpcServerControllerOpts,
			grpcx.ServerControllerWithUnixSocket(unixSocket),
		)
	}

	grpcServerController := grpcx.NewServerController(
		grpcServerAddr,
		grpcServer,
		log,
		grpcServerControllerOpts...,
	)

	if err := grpcServerController.Start(ctx); err != nil {
		log.Error("gRPC server start error", zap.Error(err))
		return exitCodeStartFailure
	}

	grpcGateway := grpcgateway.NewGateway(
		log,
		grpcgateway.WithRegistrableServices(
			"/api/v1alpha1",
			projectServiceServer,
			recordServiceServer,
		),
	)

	grpcGatewayUpstreamAddr := grpcServerAddr
	if unixSocketAvailable {
		grpcGatewayUpstreamAddr = "unix://" + unixSocket
	}

	if err := grpcGateway.ConnectAndRegister(grpcGatewayUpstreamAddr); err != nil {
		log.Error("gRPC gateway failed to connect to gRPC server", zap.Error(err))
		return exitCodeStartFailure
	}

	httpServerAddr := ":" + conf.HTTP.Port
	httpserver := httpx.NewServer(httpServerAddr, grpcGateway.Handler(), log)

	httpServerController := httpx.NewServerController(httpserver, log)
	httpServerController.Start()

	// --- Running phase ---

	slog.Infof("%s %s is started and running", appName, commandNameServer)

	sigCtx, sigCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer sigCancel()

	exitCode := exitCodeOK

	select {
	case <-sigCtx.Done():
		log.Info("Received shutdown signal")
	case err := <-grpcServerController.Wait():
		log.Error("gRPC server unexpected error", zap.Error(err))
		exitCode = exitCodeRunFailure
	case err := <-httpServerController.Wait():
		log.Error("HTTP server unexpected error", zap.Error(err))
		exitCode = exitCodeRunFailure
	}

	// --- Shutdown phase ---

	slog.Infof("%s %s is stopping...", appName, commandNameServer)

	if err := httpServerController.Stop(ctx); err != nil {
		log.Error("HTTP Server stop error", zap.Error(err))
		exitCode = exitCodeRunFailure
	}

	if err := grpcServerController.Stop(ctx); err != nil {
		log.Error("gRPC server stop error", zap.Error(err))
		exitCode = exitCodeRunFailure
	}

	if err := tracingProvider.Close(ctx); err != nil {
		log.Error("Tracing provider close error", zap.Error(err))
		exitCode = exitCodeRunFailure
	}

	if unixSocket != "" {
		if err := uds.CleanupSocket(unixSocket); err != nil {
			log.Warn("Cleanup unix socket error", zap.Error(err))
		}
	}

	slog.Infof("%s %s is stopped", appName, commandNameServer)

	return exitCode
}

func initTracing(conf TracingConfig, log *zap.Logger) (*otelx.Provider, error) {
	otelx.SetupErrorHandler(log)

	if !conf.Enabled {
		return otelx.NoopProvider(), nil
	}

	var opts []otelx.ProviderOption
	switch conf.Exporter {
	case tracingExporterLog:
		opts = append(opts, otelx.ProviderWithLogExporter(conf.Log.Pretty))
	case tracingExporterJaeger:
		opts = append(opts, otelx.ProviderWithJaegerExporter(conf.Jaeger.Endpoint))
	case tracingExporterOTLP:
		opts = append(opts, otelx.ProviderWithOTLPExporter(conf.OTLP.Endpoint))
	}

	return otelx.NewProvider(opts...)
}
