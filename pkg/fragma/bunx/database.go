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

package bunx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunotel"
	"github.com/uptrace/bun/schema"
	"go.uber.org/zap"

	"github.com/auditumio/auditum/pkg/fragma/zapx/zapxbun"
)

func NewDatabase(
	ctx context.Context,
	sqldb *sql.DB,
	dialect schema.Dialect,
	log *zap.Logger,
	logQueries LogQueriesFlag,
) (*bun.DB, error) {
	db := bun.NewDB(
		sqldb,
		dialect,
		bun.WithDiscardUnknownColumns(),
	)

	db.AddQueryHook(bunotel.NewQueryHook())
	if logQueries == LogQueriesEnabled {
		db.AddQueryHook(zapxbun.NewLogQueryHook(log))
	}

	bun.SetLogger(zapxbun.NewLogger(log))

	if err := ping(ctx, db); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}

	return db, nil
}

type LogQueriesFlag int

const (
	LogQueriesDisabled LogQueriesFlag = iota
	LogQueriesEnabled
)

func LogQueriesFlagFromBool(enabled bool) LogQueriesFlag {
	if enabled {
		return LogQueriesEnabled
	}

	return LogQueriesDisabled
}

func ping(ctx context.Context, db *bun.DB) error {
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	const maxAttempts = 100

	var err error
	for i := 0; i < maxAttempts; i++ {
		err = db.PingContext(ctx)
		if err == nil {
			return nil
		}

		time.Sleep(timeout / maxAttempts)
	}

	return err
}
