package oauth2

import (
	"net/url"

	"github.com/gorilla/schema"
)

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeClientCredentials GrantType = "client_credentials"
	GrantTypePassword          GrantType = "password"
	GrantTypeRefreshToken      GrantType = "refresh_token"

	GrantTypeUnknown GrantType = ""
)

func GrantTypeFromString(s string) GrantType {
	switch s {
	case string(GrantTypeAuthorizationCode):
		return GrantTypeAuthorizationCode
	case string(GrantTypeClientCredentials):
		return GrantTypeClientCredentials
	case string(GrantTypePassword):
		return GrantTypePassword
	case string(GrantTypeRefreshToken):
		return GrantTypeRefreshToken
	}
	return GrantTypeUnknown
}

type ErrorCode string

const (
	ErrorServerError          ErrorCode = "server_error"
	ErrorInvalidRequest       ErrorCode = "invalid_request"
	ErrorInvalidClient        ErrorCode = "invalid_client"
	ErrorInvalidGrant         ErrorCode = "invalid_grant"
	ErrorUnauthorizedClient   ErrorCode = "unauthorized_client"
	ErrorUnsupportedGrantType ErrorCode = "unsupported_grant_type"
)

type ResponseType string

const (
	ResponseTypeCode  ResponseType = "code"
	ResponseTypeToken ResponseType = "token"

	ResponseTypeUnknown ResponseType = ""
)

func ResponseTypeFromString(s string) ResponseType {
	switch s {
	case string(ResponseTypeCode):
		return ResponseTypeCode
	case string(ResponseTypeToken):
		return ResponseTypeToken
	}
	return ResponseTypeUnknown
}

func (responseType ResponseType) String() string { return string(responseType) }

type TokenType string

//NOTE: token types are case-insensitive
const (
	TokenTypeBearer TokenType = "bearer"
)

// TokenResponse is used on successful authorization. The authorization
// server issues an access token and optional refresh
// token, and constructs the response by adding the following parameters
// to the entity-body of the HTTP response with a 200 (OK) status code
type TokenResponse struct {
	// The access token issued by the authorization server.
	AccessToken string `json:"access_token" schema:"access_token"`
	// The type of the token issued as described in
	// Section 7.1.  Value is case insensitive.
	TokenType TokenType `json:"token_type" schema:"token_type"`
	// The lifetime in seconds of the access token.  For
	// example, the value "3600" denotes that the access token will
	// expire in one hour from the time the response was generated.
	// If omitted, the authorization server SHOULD provide the
	// expiration time via other means or document the default value.
	ExpiresIn int64 `json:"expires_in,omitempty" schema:"expires_in,omitempty"`
	// The refresh token, which can be used to obtain new
	// access tokens using the same authorization grant as described
	// in Section 6.
	RefreshToken string `json:"refresh_token,omitempty" schema:"refresh_token,omitempty"`
	// The scope of the access token as described by Section 3.3.
	Scope string `json:"scope,omitempty" schema:"scope,omitempty"`

	State string `json:"-" schema:"state,omitempty"`
}

func (TokenResponse) SwaggerDoc() map[string]string {
	return map[string]string{
		"": "See https://tools.ietf.org/html/rfc6749#section-5.1 for details.",
	}
}

type ErrorResponse struct {
	Error            ErrorCode `json:"error" schema:"error"`
	ErrorDescription string    `json:"error_description,omitempty" schema:"error_description,omitempty"`
	ErrorURI         string    `json:"error_uri,omitempty" schema:"error_uri,omitempty"`
	State            string    `json:"-" schema:"state,omitempty"`
}

func (ErrorResponse) SwaggerDoc() map[string]string {
	return map[string]string{
		"": "See https://tools.ietf.org/html/rfc6749#section-5.2 for details.",
	}
}

type AuthorizationRequest struct {
	ResponseType string `schema:"response_type"`
	ClientID     string `schema:"client_id"`
	RedirectURI  string `schema:"redirect_uri,omitempty"`
	Scope        string `schema:"scope,omitepmty"`
	State        string `schema:"state,omitempty"`
}

func AuthorizationRequestFromURLValues(values url.Values) (*AuthorizationRequest, error) {
	var req AuthorizationRequest
	err := schemaDecoder.Decode(&req, values)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

type AuthorizationResponse struct {
	Code  string `schema:"code"`
	State string `schema:"state"`
}

type AccessTokenRequest struct {
	GrantType GrantType `schema:"grant_type"`
	// Code is required in 'code' flow
	Code string `schema:"code"`
	// Username is required in 'password' flow
	Username string `schema:"username"`
	// Password is used in 'password' flow
	Password string `schema:"password"`
}

func QueryString(d interface{}) (queryString string, err error) {
	values := url.Values{}
	err = schemaEncoder.Encode(d, values)
	if err != nil {
		return "", err
	}
	return values.Encode(), nil
}

func MustQueryString(d interface{}) string {
	s, err := QueryString(d)
	if err != nil {
		panic(err)
	}
	return s
}

var (
	schemaEncoder = schema.NewEncoder()
	schemaDecoder = schema.NewDecoder()
)
