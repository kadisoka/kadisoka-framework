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
	NameDefault                    = "Kadisoka"
	URLDefault                     = "https://github.com/kadisoka"
	EmailDefault                   = "nop@example.com"
	NotificationEmailSenderDefault = "no-reply@example.com"
	TeamNameDefault                = "Team Kadisoka"
)

func DefaultInfo() Info {
	return Info{
		Name:                    NameDefault,
		URL:                     URLDefault,
		Email:                   EmailDefault,
		NotificationEmailSender: NotificationEmailSenderDefault,
		TeamName:                TeamNameDefault,
	}
}

type Info struct {
	// Name of the realm
	Name string
	// Canonical URL of the realm
	URL                     string
	TermsOfServiceURL       string
	PrivacyPolicyURL        string
	Email                   string
	NotificationEmailSender string
	TeamName                string
}

func InfoFromEnvOrDefault() (Info, error) {
	info := DefaultInfo()
	err := stev.LoadEnv(EnvPrefixDefault, &info)
	if err != nil {
		return DefaultInfo(), errors.Wrap("info loading from environment variables", err)
	}
	return info, nil
}
