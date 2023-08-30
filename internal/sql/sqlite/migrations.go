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

package sqlite

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file" // init driver for fs
	"github.com/uptrace/bun"

	"github.com/auditumio/auditum/pkg/fragma/bunx"
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
