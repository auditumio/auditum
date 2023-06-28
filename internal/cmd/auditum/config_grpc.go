package auditum

import (
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type GRPCConfig struct {
	Port string `yaml:"port" json:"port"`
}

func (c GRPCConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Port, validation.Required, is.Port),
	)
}

var defaultGRPCConfig = GRPCConfig{
	Port: "9090",
}
