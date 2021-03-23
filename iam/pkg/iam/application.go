package iam

import (
	"crypto/rand"
	"encoding/binary"
)

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
	ID   ApplicationRefKey
	Data ApplicationData
}

type ApplicationDataProvider interface {
	GetApplication(id ApplicationRefKey) (*Application, error)
}

// GenerateApplicationRefKey generates a new ApplicationRefKey. Note that this function is
// not consulting any database. To ensure that the generated ApplicationRefKey is
// unique, check the client database.
func GenerateApplicationRefKey(firstParty bool, clientTyp string) ApplicationRefKey {
	var typeInfo uint32
	if firstParty {
		typeInfo = _ApplicationIDNumFirstPartyBits
	}
	switch clientTyp {
	case "service":
		typeInfo |= _ApplicationIDNumServiceBits
	case "ua-public":
		typeInfo |= _ApplicationIDNumUserAgentAuthorizationPublicBits
	case "ua-confidential":
		typeInfo |= _ApplicationIDNumUserAgentAuthorizationConfidentialBits
	default:
		panic("Unsupported client app type")
	}
	instIDBytes := make([]byte, 4)
	_, err := rand.Read(instIDBytes[1:])
	if err != nil {
		panic(err)
	}
	//TODO: reserve some ranges (?)
	instID := binary.BigEndian.Uint32(instIDBytes) & ApplicationIDNumSignificantBitsMask
	appID := ApplicationIDNumFromPrimitiveValue(int32(typeInfo | instID))
	return NewApplicationRefKey(appID)
}
