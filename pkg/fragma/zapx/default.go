package zapx

import (
	"fmt"

	"go.uber.org/zap"
)

// NewDefault returns a new zap logger with sane defaults.
func NewDefault() *zap.Logger {
	conf := defaultConfig()

	log, err := conf.Build()
	if err != nil {
		panic(fmt.Sprintf("build logger: %v", err))
	}

	return log
}
