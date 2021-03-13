// package realm provides realm-related functionalities.
//
// A realm is a kind of a domain or a site. In general, it has the same
// definition as HTTP's definition of realm.
//
// A realm is not to be confused with an app github.com/kadisoka/kadisoka-framework/foundation/pkg/app .
// A realm could be comprised of various services, and each service might be
// served by many instances of the service. Each of this instance is an app.
//
// Although currently not supported, there might be a chance that we will
// also support for multiple realm in a single app.
package realm

import (
	"github.com/rez-go/stev"

	"github.com/alloyzeus/go-azfl/azfl/errors"
)

const EnvVarsPrefixDefault = "REALM_"

const (
	NameDefault                            = "UNKNOWN Realm"
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

// Info holds information about a realm.
type Info struct {
	// Name of the realm
	Name string
	// Description contains a brief textual information about the realm.
	Description string
	// Canonical URL of the realm. This URL usually refers to the website
	// of the realm.
	URL string

	TermsOfServiceURL string
	PrivacyPolicyURL  string

	ContactInfo                     ContactInfo
	DeveloperInfo                   DeveloperInfo
	NotificationEmailsSenderAddress string
}

// InfoZero returns a zero-valued Info.
func InfoZero() Info { return Info{} }

func InfoFromEnvOrDefault(envVarsPrefix string) (Info, error) {
	defaultInfo := DefaultInfo()
	return InfoFromEnv(envVarsPrefix, &defaultInfo)
}

// InfoFromEnv creates an instance of Info with values looked up from their
// respective environment variables which names are prefixed with the value
// of envVarsPrefix. The argument defaultInfo serves as the the default, that's
// it, a value in defaultInfo will be preserved unless it's overriden
// by the value specified in environment variables.
func InfoFromEnv(envVarsPrefix string, defaultInfo *Info) (Info, error) {
	if envVarsPrefix == "" {
		envVarsPrefix = EnvVarsPrefixDefault
	}
	if defaultInfo == nil {
		info := InfoZero()
		defaultInfo = &info
	}
	err := stev.LoadEnv(envVarsPrefix, defaultInfo)
	if err != nil {
		return InfoZero(), errors.Wrap("info loading from environment variables", err)
	}
	return *defaultInfo, nil
}

//TODO: person or organization.
type ContactInfo struct {
	EmailAddress string
}

//TODO: person or organization.
type DeveloperInfo struct {
	Name string
}

//TODO: person or organization.
type MaintainerInfo struct {
	Name string
}
