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

var _ azcore.EntityData = ApplicationData{}

func (cl ApplicationData) HasOAuth2RedirectURI(redirectURI string) bool {
	if cl.OAuth2RedirectURI == nil {
		return false
	}
	for _, v := range cl.OAuth2RedirectURI {
		if v == redirectURI {
			return true
		}
	}
	return false
}

type Application azcore.EntityEnvelope[ApplicationIDNum, ApplicationID, ApplicationData]

type ApplicationDataProvider interface {
	GetApplication(id ApplicationID) (*Application, error)
}
