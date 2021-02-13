// package realm provides realm-related functionalities.
//
// A realm is a kind of a domain. In general, it has the same definition
// as HTTP's definition of realm.
package realm

import (
	"github.com/rez-go/stev"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/errors"
)

const EnvPrefixDefault = "REALM_"

const (
	NameDefault                            = "Kadisoka"
	URLDefault                             = "https://github.com/kadisoka"
	ContactEmailAddressDefault             = "nop@example.com"
	NotificationEmailsSenderAddressDefault = "no-reply@example.com"
	DeveloperNameDefault                   = "Team Kadisoka"
)

func DefaultInfo() Info {
	return Info{
		Name:                            NameDefault,
		URL:                             URLDefault,
		ContactInfo:                     ContactInfo{EmailAddress: ContactEmailAddressDefault},
		DeveloperInfo:                   DeveloperInfo{Name: DeveloperNameDefault},
		NotificationEmailsSenderAddress: NotificationEmailsSenderAddressDefault,
	}
}

type Info struct {
	// Name of the realm
	Name string
	// Canonical URL of the realm
	URL                             string
	TermsOfServiceURL               string
	PrivacyPolicyURL                string
	ContactInfo                     ContactInfo
	DeveloperInfo                   DeveloperInfo
	NotificationEmailsSenderAddress string
}

func InfoFromEnvOrDefault() (Info, error) {
	info := DefaultInfo()
	err := stev.LoadEnv(EnvPrefixDefault, &info)
	if err != nil {
		return DefaultInfo(), errors.Wrap("info loading from environment variables", err)
	}
	return info, nil
}

type ContactInfo struct {
	EmailAddress string
}

type DeveloperInfo struct {
	Name string
}
