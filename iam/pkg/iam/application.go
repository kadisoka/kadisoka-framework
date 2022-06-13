package iam

import (
	"github.com/alloyzeus/go-azfl/azcore"
)

type ApplicationData struct {
	DisplayName       string
	Secret            string
	PlatformType      string // only for user-agent types
	RequiredScopes    []string
	OAuth2RedirectURI []string
}

var _ azcore.EntityAttributes = ApplicationData{}
var _ azcore.ValueObjectAssert[ApplicationData] = ApplicationData{}

func (ApplicationData) AZAttributes()       {}
func (ApplicationData) AZEntityAttributes() {}

func (appData ApplicationData) Clone() ApplicationData {
	// appData is already a clone, but it's shallow. here we are doing
	// additional operations to copy over values with shared underlying
	// instances like slices and maps
	if src := appData.RequiredScopes; src != nil {
		dst := make([]string, len(src))
		copy(dst, src)
		appData.RequiredScopes = dst
	}
	if src := appData.OAuth2RedirectURI; src != nil {
		dst := make([]string, len(src))
		copy(dst, src)
		appData.OAuth2RedirectURI = dst
	}
	return appData
}

func (appData ApplicationData) HasOAuth2RedirectURI(redirectURI string) bool {
	if appData.OAuth2RedirectURI == nil {
		return false
	}
	for _, v := range appData.OAuth2RedirectURI {
		if v == redirectURI {
			return true
		}
	}
	return false
}

type Application azcore.KeyedEntityAttributes[
	ApplicationIDNum, ApplicationID, ApplicationData]

type ApplicationDataProvider interface {
	GetApplication(id ApplicationID) (*Application, error)
}
