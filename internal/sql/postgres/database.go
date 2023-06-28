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

	"github.com/infragmo/auditum/pkg/fragma/bunx"
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

	return bunx.NewDatabase(ctx, stdlib.OpenDB(*conf), pgdialect.New(), log)
}
