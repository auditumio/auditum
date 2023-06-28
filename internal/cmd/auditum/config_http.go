package auditum

import (
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type HTTPConfig struct {
	Port string `yaml:"port" json:"port"`
}

func (c HTTPConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Port, validation.Required, is.Port),
	)
}

var defaultHTTPConfig = HTTPConfig{
	Port: "8080",
}
