package eav10n

import (
	"path/filepath"
	"reflect"
	"time"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/rez-go/stev"
)

var ResourcesDirDefault string

func init() {
	type t int
	pkgPath := reflect.TypeOf(t(0)).PkgPath()
	ResourcesDirDefault = filepath.Join(pkgPath, "resources")
}

func ConfigFromEnv(prefix string, seedCfg *Config) (*Config, error) {
	if seedCfg == nil {
		seedCfg = &Config{}
	}
	err := stev.LoadEnv(prefix, seedCfg)
	if err != nil {
		return nil, errors.Wrap("config loading from environment variables", err)
	}
	return seedCfg, nil
}

type Config struct {
	CodeTTLDefault          time.Duration `env:"CODE_TTL_DEFAULT"`
	ConfirmationAttemptsMax int16         `env:"CONFIRMATION_ATTEMPTS_MAX,docs_hidden"`
	SenderAddress           string        `env:"SENDER_ADDRESS"`
	ResourcesDir            string        `env:"RESOURCES_DIR"`
	// The email delivery service to use.
	EmailDeliveryService string `env:"EMAIL_DELIVERY_SERVICE,required"`
	// Configurations for modules
	Modules map[string]ModuleConfig `env:",map,squash"`
}

func (cfg Config) FieldDocsDescriptor(fieldName string) *stev.FieldDocsDescriptor {
	switch fieldName {
	case "EmailDeliveryService", "EMAIL_DELIVERY_SERVICE":
		modules := map[string]stev.EnumValueDocs{}
		for k, v := range cfg.Modules {
			var shortDesc string
			if smsCfg := v.EmailDeliveryServiceConfig(); smsCfg != nil {
				if moduleDesc := stev.LoadSelfDocsDescriptor(smsCfg); moduleDesc != nil {
					shortDesc = moduleDesc.ShortDesc
				}
			}
			modules[k] = stev.EnumValueDocs{ShortDesc: shortDesc}
		}
		return &stev.FieldDocsDescriptor{
			Description:     "The email delivery service to use.",
			AvailableValues: modules,
		}
	}
	return nil
}
