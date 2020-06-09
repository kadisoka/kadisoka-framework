package rest

import (
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/rez-go/stev"

	"github.com/citadelium/foundation/pkg/errors"
)

type CORSFilterConfig struct {
	AllowedHeaders *string `env:"ALLOWED_HEADERS"`
	AllowedMethods string  `env:"ALLOWED_METHODS"`
	AllowedDomains string  `env:"ALLOWED_DOMAINS"`
}

func SetUpCORSFilterByEnv(restContainer *restful.Container, envPrefix string) error {
	var cfg CORSFilterConfig
	err := stev.LoadEnv(envPrefix, &cfg)
	if err != nil {
		return errors.Wrap("config loading from environment variables", err)
	}

	var allowedHeaders []string
	if cfg.AllowedHeaders != nil {
		if strVal := *cfg.AllowedHeaders; strVal != "" {
			parts := strings.Split(strVal, ",")
			for _, str := range parts {
				allowedHeaders = append(allowedHeaders, strings.TrimSpace(str))
			}
		}
	} else {
		// These are what we generally need
		allowedHeaders = []string{"Content-Type", "Accept", "Authorization"}
	}

	var allowedMethods []string
	if strVal := cfg.AllowedMethods; strVal != "" {
		parts := strings.Split(strVal, ",")
		for _, str := range parts {
			allowedMethods = append(allowedMethods, strings.TrimSpace(str))
		}
	}

	var allowedDomains []string
	if strVal := cfg.AllowedDomains; strVal != "" {
		parts := strings.Split(strVal, ",")
		for _, str := range parts {
			allowedDomains = append(allowedDomains, strings.TrimSpace(str))
		}
	}

	restContainer.Filter(restful.CrossOriginResourceSharing{
		AllowedHeaders: allowedHeaders,
		AllowedMethods: allowedMethods,
		AllowedDomains: allowedDomains,
		CookiesAllowed: false,
		Container:      restContainer}.Filter)

	return nil
}
