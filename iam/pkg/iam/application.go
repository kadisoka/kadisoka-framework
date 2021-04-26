package iam

type ApplicationData struct {
	DisplayName       string
	Secret            string
	PlatformType      string // only for user-agent types
	RequiredScopes    []string
	OAuth2RedirectURI []string
}

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

type Application struct {
	RefKey ApplicationRefKey
	Data   ApplicationData
}

type ApplicationDataProvider interface {
	GetApplication(refKey ApplicationRefKey) (*Application, error)
}
