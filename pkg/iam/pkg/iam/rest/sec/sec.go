package sec

type AuthorizationType string

func (opSecType AuthorizationType) String() string {
	return string(opSecType)
}

const (
	AuthorizationBasicOAuth2ClientCredentials AuthorizationType = "basic-oauth2-client-credentials"
	AuthorizationBearerAccessToken            AuthorizationType = "bearer-access-token"
)
