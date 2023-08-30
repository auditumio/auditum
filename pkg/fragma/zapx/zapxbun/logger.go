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

package zapxbun

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(log *zap.Logger) *Logger {
	return &Logger{
		log: log.Named("bun").Sugar(),
	}
}

// Logger implements github.com/uptrace/bun logger.
type Logger struct {
	log *zap.SugaredLogger
}

func (a *Logger) Printf(format string, args ...any) {
	a.log.Infof(format, args...)
}

func NewLogQueryHook(log *zap.Logger) bun.QueryHook {
	return &LogQueryHook{log: log.Named("bun_sql")}
}

// LogQueryHook implements bun query hook that logs queries.
type LogQueryHook struct {
	log *zap.Logger
}

func (l *LogQueryHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (l *LogQueryHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	fields := []zapcore.Field{
		zap.String("query", event.Query),
		zap.Int64("duration_ms", time.Since(event.StartTime).Milliseconds()),
	}

	if event.Result != nil {
		if ra, err := event.Result.RowsAffected(); err == nil {
			fields = append(fields, zap.Int64("rows_affected", ra))
		}
	}

	level := zapcore.InfoLevel
	if event.Err != nil &&
		!errors.Is(event.Err, sql.ErrNoRows) {
		level = zapcore.ErrorLevel
		fields = append(fields, zap.Error(event.Err))
	}

	l.log.Log(level, "SQL query", fields...)
}
