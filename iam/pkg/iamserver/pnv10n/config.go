package pnv10n

import (
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/rez-go/stev"
)

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
	// The default code TTL
	CodeTTLDefault time.Duration `env:"CODE_TTL_DEFAULT"`
	// The maximum number of failed attempts for a verification request
	ConfirmationAttemptsMax int16 `env:"CONFIRMATION_ATTEMPTS_MAX,docs_hidden"`
	// For use with SMS Retriever API https://developers.google.com/identity/sms-retriever/overview
	SMSRetrieverAppHash string `env:"SMS_RETRIEVER_APP_HASH"`
	// The SMS delivery service to use.
	SMSDeliveryService string `env:"SMS_DELIVERY_SERVICE,required"`
	// Configurations for modules
	Modules map[string]interface{} `env:",map,squash"`
}

func (cfg Config) FieldDocsDescriptor(fieldName string) *stev.FieldDocsDescriptor {
	switch fieldName {
	case "SMSDeliveryService", "SMS_DELIVERY_SERVICE":
		modules := map[string]stev.EnumValueDocs{}
		for k := range cfg.Modules {
			modules[k] = stev.EnumValueDocs{}
		}
		return &stev.FieldDocsDescriptor{
			Description:     "The SMS delivery service to use.",
			AvailableValues: modules,
		}
	}
	return nil
}
