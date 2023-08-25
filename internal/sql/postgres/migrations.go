package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // init driver for fs
	"github.com/uptrace/bun"

	"github.com/auditumio/auditum/pkg/fragma/bunx"
)

func RunMigrations(db *bun.DB, migrationsDir string, log any) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir, // e.g. "file://./postgres/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance: %v", err)
	}

	return bunx.RunMigrations(m, log)
}
