//go:build integration && postgres && !sqlite

package sqltest

import (
	"context"
	"sync"
	"testing"

	"github.com/caarlos0/env/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"go.uber.org/zap"

	"github.com/infragmo/auditum/internal/sql/postgres"
)

type configuration struct {
	Postgres postgresConfig `envPrefix:"POSTGRES_"`
}

type postgresConfig struct {
	Host     string `env:"HOST" envDefault:"127.0.0.1"`
	Port     string `env:"PORT" envDefault:"5432"`
	Database string `env:"DATABASE" envDefault:"auditum_db"`
	Username string `env:"USERNAME" envDefault:"user"`
	Password string `env:"PASSWORD" envDefault:"pass"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`

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

	db, err := postgres.NewDatabase(ctx,
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.Database,
		conf.Postgres.Username,
		conf.Postgres.Password,
		conf.Postgres.SSLMode,
		zap.NewNop(),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		assert.NoError(t, err)
	})

	if conf.Postgres.LogQueries {
		db.AddQueryHook(logQueriesQueryHook{t: t})
	}

	// Run migrations.

	migrationsOnce.Do(func() {
		runMigrations(t, db, conf.Postgres.LogQueries)
	})

	return db
}

func runMigrations(t *testing.T, db *bun.DB, logQueries bool) {
	t.Helper()

	var log any
	if logQueries {
		log = migrateLogger{t: t}
	}

	err := postgres.RunMigrations(db, "./postgres/migrations", log)
	require.NoError(t, err)
}
