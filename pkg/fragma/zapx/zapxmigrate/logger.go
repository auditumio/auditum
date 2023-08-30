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

package zapxmigrate

import "go.uber.org/zap"

func NewLogger(log *zap.Logger) *Logger {
	return &Logger{
		log: log.Named("migrate").Sugar(),
	}
}

// Logger implements github.com/golang-migrate/migrate/v4 logger.
type Logger struct {
	log *zap.SugaredLogger
}

func (a *Logger) Printf(format string, v ...interface{}) {
	a.log.Infof(format, v...)
}

func (a *Logger) Verbose() bool {
	return true
}
