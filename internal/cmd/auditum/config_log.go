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
