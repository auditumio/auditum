//go:build integration

package sqltest

import (
	"context"
	"testing"

	"github.com/uptrace/bun"
)

type logQueriesQueryHook struct {
	t *testing.T
}

func (l logQueriesQueryHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (l logQueriesQueryHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	l.t.Log(event.Query)
}

type migrateLogger struct {
	t *testing.T
}

func (m migrateLogger) Printf(format string, v ...interface{}) {
	m.t.Logf(format, v...)
}

func (m migrateLogger) Verbose() bool {
	return true
}
