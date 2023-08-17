package auditum

import (
	"github.com/invopop/validation"
	"go.uber.org/zap"

	"github.com/auditumio/auditum/pkg/fragma/zapx"
)

type LogConfig struct {
	Format string `yaml:"format" json:"format"`
	Level  string `yaml:"level" json:"level"`
}

func (c LogConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Format,
			validation.Required,
			validation.In(
				zapx.FormatJSON,
				zapx.FormatText,
			),
		),
		validation.Field(
			&c.Level,
			validation.Required,
			validation.In(
				zap.DebugLevel.String(),
				zap.InfoLevel.String(),
				zap.WarnLevel.String(),
				zap.ErrorLevel.String(),
			),
		),
	)
}

var defaultLogConfig = LogConfig{
	Format: zapx.FormatJSON,
	Level:  zap.InfoLevel.String(),
}
