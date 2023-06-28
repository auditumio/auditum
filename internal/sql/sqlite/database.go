package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"

	"github.com/infragmo/auditum/pkg/fragma/bunx"
)

const FilepathMemory = ":memory:"

func NewDatabase(
	ctx context.Context,
	fpath string, // e.g. :memory: or /tmp/auditum.db
	log *zap.Logger,
) (*bun.DB, error) {
	dsn := fmt.Sprintf("file:%s?cache=shared", fpath)

	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open sql db: %v", err)
	}

	// Required for in-memory database.
	sqldb.SetMaxIdleConns(1000)
	sqldb.SetConnMaxLifetime(0)

	return bunx.NewDatabase(ctx, sqldb, sqlitedialect.New(), log)
}
