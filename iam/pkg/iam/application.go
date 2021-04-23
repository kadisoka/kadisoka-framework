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

// GenerateApplicationRefKey generates a new ApplicationRefKey. Note that this function is
// not consulting any database. To ensure that the generated ApplicationRefKey is
// unique, check the client database.
func GenerateApplicationRefKey(firstParty bool, clientTyp string) ApplicationRefKey {
	var typeInfo uint32
	if firstParty {
		typeInfo = ApplicationIDNumFirstPartyBits
	}
	switch clientTyp {
	case "service":
		typeInfo |= ApplicationIDNumServiceBits
	case "ua-public":
		typeInfo |= ApplicationIDNumUserAgentAuthorizationPublicBits
	case "ua-confidential":
		typeInfo |= ApplicationIDNumUserAgentAuthorizationConfidentialBits
	default:
		panic("Unsupported client app type")
	}
	//TODO: reserve some ranges (?)
	appIDNum, err := GenerateApplicationIDNum(typeInfo)
	if err != nil {
		panic(err)
	}
	return NewApplicationRefKey(appIDNum)
}
