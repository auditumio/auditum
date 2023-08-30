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
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"

	"github.com/auditumio/auditum/pkg/fragma/bunx"
)

const FilepathMemory = ":memory:"

func NewDatabase(
	ctx context.Context,
	fpath string, // e.g. :memory: or /tmp/auditum.db
	log *zap.Logger,
	logQueries bunx.LogQueriesFlag,
) (*bun.DB, error) {
	dsn := fmt.Sprintf("file:%s?cache=shared", fpath)

	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open sql db: %v", err)
	}

	// Required for in-memory database.
	sqldb.SetMaxIdleConns(1000)
	sqldb.SetConnMaxLifetime(0)

	return bunx.NewDatabase(
		ctx,
		sqldb,
		sqlitedialect.New(),
		log,
		logQueries,
	)
}
