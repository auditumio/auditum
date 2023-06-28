package sqlite

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file" // init driver for fs
	"github.com/uptrace/bun"

	"github.com/infragmo/auditum/pkg/fragma/bunx"
)

func RunMigrations(db *bun.DB, fpath string, migrationsDir string, log any) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{
		DatabaseName: fpath,
		NoTxWrap:     true,
	})
	if err != nil {
		return fmt.Errorf("create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir, // e.g. "file://./sqlite/migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance: %v", err)
	}

	return bunx.RunMigrations(m, log)
}
