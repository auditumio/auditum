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

package auditum

import (
	"go.uber.org/zap"

	"github.com/auditumio/auditum/pkg/fragma/zapx"
)

func Execute() int {
	dlog := zapx.NewDefault()
	defer func() {
		_ = dlog.Sync()
	}()

	cmd, configPath, err := parseCommand()
	if err != nil {
		dlog.Error("Failed to parse command", zap.Error(err))
		return exitCodeStartFailure
	}

	conf, err := loadConfiguration(configPath)
	if err != nil {
		dlog.Error("Failed to load configuration", zap.Error(err))
		return exitCodeStartFailure
	}

	log, err := zapx.New(
		zapx.WithFormat(conf.Log.Format),
		zapx.WithLevel(conf.Log.Level),
	)
	if err != nil {
		dlog.Error("Failed to create logger", zap.Error(err))
		return exitCodeStartFailure
	}
	defer func() {
		_ = log.Sync()
	}()

	return executeCommand(cmd, conf, log)
}
