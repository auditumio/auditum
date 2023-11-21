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

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

const (
	tracingExporterLog    = "log"
	tracingExporterJaeger = "jaeger"
	tracingExporterOTLP   = "otlp"
)

type TracingConfig struct {
	Enabled  bool                        `yaml:"enabled" json:"enabled"`
	Exporter string                      `yaml:"exporter" json:"exporter"`
	Log      TracingLogExporterConfig    `yaml:"stdout" json:"stdout"`
	Jaeger   TracingJaegerExporterConfig `yaml:"jaeger" json:"jaeger"`
	OTLP     TracingOTLPExporterConfig   `yaml:"otlp" json:"otlp"`
}

func (c TracingConfig) Validate() error {
	err := validation.ValidateStruct(&c,
		validation.Field(
			&c.Exporter,
			validation.Required,
			validation.In(
				tracingExporterLog,
				tracingExporterJaeger,
				tracingExporterOTLP,
			),
		),
	)
	if err != nil {
		return err
	}

	switch c.Exporter {
	case tracingExporterLog:
		// No additional validation required.
		return nil
	case tracingExporterJaeger:
		if err := c.Jaeger.Validate(); err != nil {
			return fmt.Errorf("invalid 'jaeger': %v", err)
		}
		return nil
	case tracingExporterOTLP:
		if err := c.OTLP.Validate(); err != nil {
			return fmt.Errorf("invalid 'otlp': %v", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown 'exporter': %s", c.Exporter)
	}
}

type TracingLogExporterConfig struct {
	Pretty bool `yaml:"pretty" json:"pretty"`
}

type TracingJaegerExporterConfig struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

func (c TracingJaegerExporterConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Endpoint, validation.Required, is.URL),
	)
}

type TracingOTLPExporterConfig struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

func (c TracingOTLPExporterConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Endpoint,
			validation.Required,
			// is.URL does not work: under the hood govalidator.IsURL is used
			// which has regex which permits only a few schemes, and does not
			// include grpc.
			is.RequestURL,
		),
	)
}

var defaultTracingConfig = TracingConfig{
	Enabled:  false,
	Exporter: tracingExporterLog,
	Log: TracingLogExporterConfig{
		Pretty: false,
	},
	Jaeger: TracingJaegerExporterConfig{
		Endpoint: "http://localhost:14268/api/traces",
	},
	OTLP: TracingOTLPExporterConfig{
		Endpoint: "http://localhost:4318",
	},
}
