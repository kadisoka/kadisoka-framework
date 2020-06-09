package store

import (
	"github.com/rez-go/stev"
)

// ConfigFromEnv populate the configuration by looking up the environment variables.
func ConfigFromEnv(envPrefix string) (*Config, error) {
	cfg := Config{}
	err := stev.LoadEnv(envPrefix, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Config File Core cofiguration
type Config struct {
	ImagesBaseURL string `env:"IMAGES_BASE_URL"`
	//TODO: declare the encoding (hex?)
	NameGenerationKey string                 `env:"FILENAME_GENERATION_KEY"`
	StoreService      string                 `env:"STORE_SERVICE"`
	Modules           map[string]interface{} `env:",map,squash"`
}
