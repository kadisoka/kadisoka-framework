package rest

import (
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/rez-go/stev"

	"github.com/alloyzeus/go-azfl/azfl/errors"
)

type CORSFilterConfig struct {
	AllowedHeaders *string `env:"ALLOWED_HEADERS"`
	AllowedMethods string  `env:"ALLOWED_METHODS"`
	AllowedDomains string  `env:"ALLOWED_DOMAINS"`
}

func SetUpCORSFilter(
	restContainer *restful.Container,
	filterConfig CORSFilterConfig,
) error {
	var allowedHeaders []string
	if filterConfig.AllowedHeaders != nil {
		if strVal := *filterConfig.AllowedHeaders; strVal != "" {
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
	if strVal := filterConfig.AllowedMethods; strVal != "" {
		parts := strings.Split(strVal, ",")
		for _, str := range parts {
			allowedMethods = append(allowedMethods, strings.TrimSpace(str))
		}
	}

	var allowedDomains []string
	if strVal := filterConfig.AllowedDomains; strVal != "" {
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

func SetUpCORSFilterByEnv(
	restContainer *restful.Container,
	envVarsPrefix string,
	defaultFilterConfig *CORSFilterConfig,
) error {
	if defaultFilterConfig == nil {
		defaultFilterConfig = &CORSFilterConfig{}
	}
	err := stev.LoadEnv(envVarsPrefix, defaultFilterConfig)
	if err != nil {
		return errors.Wrap("config loading from environment variables", err)
	}

	return SetUpCORSFilter(restContainer, *defaultFilterConfig)
}
