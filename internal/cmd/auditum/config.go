package auditum

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"

	"github.com/infragmo/auditum/internal/aud"
)

const (
	envPrefix         = "AUDITUM_"
	envConfigFilepath = envPrefix + "CONFIG"
)

type Configuration struct {
	Log      LogConfig     `yaml:"log" json:"log"`
	Tracing  TracingConfig `yaml:"tracing" json:"tracing"`
	HTTP     HTTPConfig    `yaml:"http" json:"http"`
	GRPC     GRPCConfig    `yaml:"grpc" json:"grpc"`
	Store    StoreConfig   `yaml:"store" json:"store"`
	Settings aud.Settings  `yaml:"settings" json:"settings"`

	// Note: json tag in structs is used by validation package.
}

func (c Configuration) Validate() error {
	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("invalid 'log': %v", err)
	}

	if err := c.Tracing.Validate(); err != nil {
		return fmt.Errorf("invalid 'tracing': %v", err)
	}

	if err := c.HTTP.Validate(); err != nil {
		return fmt.Errorf("invalid 'http': %v", err)
	}

	if err := c.GRPC.Validate(); err != nil {
		return fmt.Errorf("invalid 'grpc': %v", err)
	}

	if err := c.Store.Validate(); err != nil {
		return fmt.Errorf("invalid 'store': %v", err)
	}

	if err := c.Settings.Validate(); err != nil {
		return fmt.Errorf("invalid 'settings': %v", err)
	}

	return nil
}

// NOTE: must be in sync with config/auditum.yaml
var defaultConfig = Configuration{
	Log:      defaultLogConfig,
	Tracing:  defaultTracingConfig,
	HTTP:     defaultHTTPConfig,
	GRPC:     defaultGRPCConfig,
	Store:    defaultStoreConfig,
	Settings: aud.DefaultSettings,
}

func loadConfiguration(fpath string) (*Configuration, error) {
	if fpath == "" {
		fpath = os.Getenv(envConfigFilepath)
	}

	// Load configuration from sources, considering default configuration.
	// Then save it into the structure.

	k := koanf.New(".")

	if err := k.Load(structs.Provider(defaultConfig, "yaml"), nil); err != nil {
		return nil, fmt.Errorf("load default configuration: %v", err)
	}
	if fpath != "" {
		if err := k.Load(file.Provider(fpath), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("load yaml file configuration: %v", err)
		}
	}
	if err := k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.ReplaceAll(
			strings.TrimPrefix(s, envPrefix),
			"_",
			".",
		)
	}), nil); err != nil {
		return nil, fmt.Errorf("load environment variables configuration: %v", err)
	}

	var conf Configuration
	if err := k.Unmarshal("", &conf); err != nil {
		return nil, fmt.Errorf("unmarshal merged configuration: %v", err)
	}

	// Validate configuration.

	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %v", err)
	}

	return &conf, nil
}
