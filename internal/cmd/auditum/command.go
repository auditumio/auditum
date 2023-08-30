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
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	appName             = "auditum"
	commandNameServer   = "server"
	commandNameMigrator = "migrator"
)

const (
	exitCodeOK = iota
	exitCodeStartFailure
	exitCodeRunFailure
)

func parseCommand() (command, configPath string, err error) {
	flagset := flag.NewFlagSet(appName, flag.ContinueOnError)

	flagset.String("config", "", "Path to config file.")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return "", "", fmt.Errorf("parse command line arguments: %v", err)
	}

	c := flagset.Arg(0)

	fpath, err := flagset.GetString("config")
	if err != nil {
		return "", "", fmt.Errorf("get config argument: %v", err)
	}

	return c, fpath, nil
}

func executeCommand(cmd string, config *Configuration, log *zap.Logger) int {
	switch cmd {
	case "server", "serve", "":
		return executeServer(config, log)
	case "migrator", "migrate":
		return executeMigrator(config, log)
	default:
		log.Error("Unknown command", zap.String("command", cmd))
		return exitCodeStartFailure
	}
}
