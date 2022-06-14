// Package eav10n provides utilities for verifying email addresses.
package eav10n

import (
	"golang.org/x/text/language"
)

const (
	messageCharset = "UTF-8"
)

var (
	messageLocaleDefault = language.MustParse("en-US")
)

const emailDeliveryServiceDefault = ""

func ConfigSkeleton() Config {
	moduleConfigs := ModuleConfigSkeletons()
	return Config{
		EmailDeliveryService: emailDeliveryServiceDefault,
		Modules:              moduleConfigs,
	}
}
