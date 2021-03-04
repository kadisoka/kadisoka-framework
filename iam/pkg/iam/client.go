package iam

type Client struct {
	ID                ApplicationRefKey
	DisplayName       string
	Secret            string
	PlatformType      string // only for user-agent types
	RequiredScopes    []string
	OAuth2RedirectURI []string
}

func (cl Client) HasOAuth2RedirectURI(redirectURI string) bool {
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

type ClientDataProvider interface {
	GetClient(id ApplicationRefKey) (*Client, error)
}
