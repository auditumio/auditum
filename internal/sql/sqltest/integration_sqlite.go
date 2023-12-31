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

//go:build integration && !postgres && sqlite

package sqltest

import (
	"context"
	"sync"
	"testing"

	"github.com/caarlos0/env/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"go.uber.org/zap"

	"github.com/auditumio/auditum/internal/sql/sqlite"
	"github.com/auditumio/auditum/pkg/fragma/bunx"
)

type configuration struct {
	SQLite sqliteConfig `envPrefix:"SQLITE_"`
}

type sqliteConfig struct {
	// e.g. ":memory:" or "auditum.db"
	Filepath string `env:"FILEPATH" envDefault:":memory:"`

	// e.g. "1" or ""
	LogQueries bool `env:"LOG_QUERIES" envDefault:"false"`
}

func loadConfiguration(t *testing.T) *configuration {
	t.Helper()

	var conf configuration

	err := env.Parse(&conf)
	require.NoError(t, err)

	return &conf
}

var (
	migrationsOnce = sync.Once{}
)

func NewDatabase(ctx context.Context, t *testing.T) *bun.DB {
	t.Helper()

	// Create database connection.

	conf := loadConfiguration(t)

	db, err := sqlite.NewDatabase(
		ctx,
		conf.SQLite.Filepath,
		zap.NewNop(),
		bunx.LogQueriesDisabled, // We add hook below.
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		assert.NoError(t, err)
	})

	if conf.SQLite.LogQueries {
		db.AddQueryHook(logQueriesQueryHook{t: t})
	}

	// Run migrations.

	if conf.SQLite.Filepath == sqlite.FilepathMemory {
		// Because in-memory database is destroyed after closing connection,
		// we need to run migrations every time.
		runMigrations(t, db, conf.SQLite.Filepath, conf.SQLite.LogQueries)
	} else {
		migrationsOnce.Do(func() {
			runMigrations(t, db, conf.SQLite.Filepath, conf.SQLite.LogQueries)
		})
	}

	return db
}

func runMigrations(t *testing.T, db *bun.DB, fpath string, logQueries bool) {
	t.Helper()

	var log any
	if logQueries {
		log = migrateLogger{t: t}
	}

	err := sqlite.RunMigrations(db, fpath, "./sqlite/migrations", log)
	require.NoError(t, err)
}
