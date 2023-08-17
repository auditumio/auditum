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
