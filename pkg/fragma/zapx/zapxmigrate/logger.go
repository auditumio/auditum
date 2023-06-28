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
