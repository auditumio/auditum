package zapxbun

import (
	"go.uber.org/zap"
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
