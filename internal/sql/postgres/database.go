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

package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/zap"

	"github.com/auditumio/auditum/pkg/fragma/bunx"
)

func NewDatabase(
	ctx context.Context,
	host string,
	port string,
	database string,
	username string,
	password string,
	sslMode string,
	log *zap.Logger,
	logQueries bunx.LogQueriesFlag,
) (*bun.DB, error) {
	dsnWords := []string{
		"host=" + host,
		"port=" + port,
		"dbname=" + database,
		"user=" + username,
		"password=" + password,
		"sslmode=" + sslMode,
		// See: https://bun.uptrace.dev/postgres/#pgx
		"default_query_exec_mode=simple_protocol",
	}

	dsn := strings.Join(dsnWords, " ")

	conf, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %v", err)
	}

	return bunx.NewDatabase(
		ctx,
		stdlib.OpenDB(*conf),
		pgdialect.New(),
		log,
		logQueries,
	)
}
