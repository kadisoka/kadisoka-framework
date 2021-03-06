// Package pnv10n provides utility for verifying phone numbers.
//
// pnv10n = phone-number verification.
package pnv10n

import (
	"golang.org/x/text/language"
)

var messageLocaleDefault = language.MustParse("en-US")

const smsDeliveryServiceDefault = ""

func ConfigSkeleton() Config {
	moduleConfigs := ModuleConfigSkeletons()
	return Config{
		SMSDeliveryService: smsDeliveryServiceDefault,
		Modules:            moduleConfigs,
	}
}
