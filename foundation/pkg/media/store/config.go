package store

import (
	"github.com/rez-go/stev"
)

// ConfigFromEnv populate the configuration by looking up the environment variables.
func ConfigFromEnv(envVarsPrefix string) (*Config, error) {
	cfg := Config{}
	err := stev.LoadEnv(envVarsPrefix, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Config File Core cofiguration
type Config struct {
	//TODO: declare the encoding (hex?)
	NameGenerationKey string `env:"FILENAME_GENERATION_KEY"`
	StoreService      string `env:"STORE_SERVICE"`

	Modules map[string]interface{} `env:",map,squash"`

	ImagesBaseURL string `env:"IMAGES_BASE_URL"`
}

func (Config) FieldDescriptions() map[string]string {
	return map[string]string{
		"NameGenerationKey": "Uploaded files are given names based on their " +
			"respective hash values to reduce duplications. This might cause " +
			"privacy issue as those who have the same file could look up it " +
			"in the server. To reduce the possible risk, we use HMAC to " +
			"generate the filename. HMAC requires key which should be " +
			"treated as a secret.\n\nThe value must be provided as a " +
			"standard base64-encoded string, padded or not.",
	}
}
