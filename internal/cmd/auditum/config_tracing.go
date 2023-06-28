package auditum

import (
	"fmt"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

const (
	tracingExporterLog    = "log"
	tracingExporterJaeger = "jaeger"
)

type TracingConfig struct {
	Enabled  bool                        `yaml:"enabled" json:"enabled"`
	Exporter string                      `yaml:"exporter" json:"exporter"`
	Log      TracingLogExporterConfig    `yaml:"stdout" json:"stdout"`
	Jaeger   TracingJaegerExporterConfig `yaml:"jaeger" json:"jaeger"`
}

func (c TracingConfig) Validate() error {
	err := validation.ValidateStruct(&c,
		validation.Field(
			&c.Exporter,
			validation.Required,
			validation.In(tracingExporterLog, tracingExporterJaeger),
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

var defaultTracingConfig = TracingConfig{
	Enabled:  false,
	Exporter: tracingExporterLog,
	Log: TracingLogExporterConfig{
		Pretty: false,
	},
	Jaeger: TracingJaegerExporterConfig{
		Endpoint: "http://localhost:14268/api/traces",
	},
}
