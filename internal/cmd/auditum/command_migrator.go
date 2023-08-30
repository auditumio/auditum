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

	"go.uber.org/zap"

	"github.com/auditumio/auditum/internal/sql/postgres"
	"github.com/auditumio/auditum/internal/sql/sqlite"
	"github.com/auditumio/auditum/pkg/fragma/bunx"
)

func executeMigrator(conf *Configuration, log *zap.Logger) (code int) {
	ctx := context.Background()

	slog := log.Sugar()
	slog.Infof("%s %s started", appName, commandNameMigrator)
	defer func() {
		if code == exitCodeOK {
			slog.Infof("%s %s finished", appName, commandNameMigrator)
		} else {
			slog.Errorf("%s %s failed", appName, commandNameMigrator)
		}
	}()

	switch conf.Store.Type {
	case storeTypeSQLite:
		return migrateSQLite(ctx, conf, log)
	case storeTypePostgres:
		return migratePostgres(ctx, conf, log)
	default:
		log.Panic("Unreachable code: invalid store type", zap.String("store_type", conf.Store.Type))
		return exitCodeStartFailure
	}
}

func migrateSQLite(ctx context.Context, conf *Configuration, log *zap.Logger) int {
	if conf.Store.SQLite.DatabasePath == sqlite.FilepathMemory {
		log.Info("Migrations for in-memory SQLite database are run automatically on server startup." +
			" The migrator command is no-op.")
		return exitCodeOK
	}

	db, err := sqlite.NewDatabase(
		ctx,
		conf.Store.SQLite.DatabasePath,
		log,
		bunx.LogQueriesFlagFromBool(conf.Store.SQLite.LogQueries),
	)
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return exitCodeStartFailure
	}

	if err := sqlite.RunMigrations(
		db,
		conf.Store.SQLite.DatabasePath,
		conf.Store.SQLite.MigrationsPath,
		log,
	); err != nil {
		log.Error("Failed to run migrations", zap.Error(err))
		return exitCodeRunFailure
	}

	return exitCodeOK
}

func migratePostgres(ctx context.Context, conf *Configuration, log *zap.Logger) int {
	db, err := postgres.NewDatabase(
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

	if err := postgres.RunMigrations(
		db,
		conf.Store.Postgres.MigrationsPath,
		log,
	); err != nil {
		log.Error("Failed to run migrations", zap.Error(err))
		return exitCodeRunFailure
	}

	return exitCodeOK
}
