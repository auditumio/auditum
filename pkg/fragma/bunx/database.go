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

	"github.com/infragmo/auditum/pkg/fragma/zapx/zapxbun"
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
