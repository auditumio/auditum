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

package zapxgrpc

import "go.uber.org/zap"

func NewLogger(log *zap.Logger) *Logger {
	return &Logger{
		slog:      log.Named("grpc_log").Sugar(),
		verbosity: 0,
	}
}

// Logger implements google.golang.org/grpc/grpclog.LoggerV2.
type Logger struct {
	slog      *zap.SugaredLogger
	verbosity int
}

func (l *Logger) Info(args ...interface{}) {
	l.slog.Info(args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.slog.Infoln(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.slog.Infof(format, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.slog.Warn(args...)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.slog.Warnln(args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.slog.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.slog.Error(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.slog.Errorln(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.slog.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.slog.Fatal(args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.slog.Fatalln(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.slog.Fatalf(format, args...)
}

func (l *Logger) V(level int) bool {
	return l.verbosity <= level
}
