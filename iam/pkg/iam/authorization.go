package iam

import (
	"time"

	"github.com/alloyzeus/go-azcore/azcore/errors"
	dataerrs "github.com/alloyzeus/go-azcore/azcore/errors/data"
	"github.com/square/go-jose/v3/jwt"
)

// Used in API call metadata: HTTP header and gRPC call metadata
const (
	AuthorizationMetadataKey    = "Authorization"
	AuthorizationMetadataKeyAlt = "authorization"
)

var (
	ErrReqFieldAuthorizationMalformed = ReqFieldErr("Authorization", dataerrs.ErrMalformed)

	ErrReqFieldAuthorizationTypeUnsupported = ReqFieldErr("Authorization", dataerrs.ErrTypeUnsupported)

	ErrAuthorizationCodeAlreadyClaimed = errors.EntMsg("authorization code", "already claimed")
)

// Authorization is generally used to provide authorization information
// for call or request. An Authorization is usually obtained from authorization
// token / access token provided along the request / call.
type Authorization struct {
	// If this context is an assumed context, this field
	// holds info about the assuming context.
	AssumingAuthorization *Authorization `json:"assuming_authorization,omitempty"`

	Session SessionRefKey

	// Scope, expiry time

	rawToken string
}

// newEmptyAuthorization creates a new instance of Authorization without
// any data.
func newEmptyAuthorization() *Authorization {
	return &Authorization{}
}

func (authCtx Authorization) IsValid() bool {
	return authCtx.Session.IsValid()
}

func (authCtx Authorization) IsNotValid() bool {
	return !authCtx.IsValid()
}

func (authCtx Authorization) Actor() Actor {
	return Actor{
		UserRef:     authCtx.Session.terminal.user,
		TerminalRef: authCtx.Session.terminal,
	}
}

// IsUserContext is used to determine if this context represents a user.
func (authCtx Authorization) IsUserContext() bool {
	if authCtx.ClientID().IsUserAgent() && authCtx.Session.terminal.user.IsValid() {
		return true
	}
	return false
}

func (authCtx Authorization) IsServiceClientContext() bool {
	if authCtx.ClientID().IsService() && authCtx.Session.terminal.user.IsNotValid() {
		return true
	}
	return false
}

func (authCtx Authorization) UserRef() UserRefKey {
	return authCtx.Session.terminal.user
}

// UserRefKeyPtr returns a pointer to a new copy of user ID. The
// returned value is non-nil when the user ref-key is valid.
func (authCtx Authorization) UserRefKeyPtr() *UserRefKey {
	if authCtx.Session.terminal.user.IsValid() {
		return &authCtx.Session.terminal.user
	}
	return nil
}

func (authCtx Authorization) UserID() UserID {
	return authCtx.Session.terminal.user.ID()
}

// UserIDPtr returns a pointer to a new copy of user ID. The
// returned value is non-nil when the user ref-key is valid.
func (authCtx Authorization) UserIDPtr() *UserID {
	return authCtx.Session.terminal.user.IDPtr()
}

func (authCtx Authorization) TerminalRef() TerminalRefKey {
	return authCtx.Session.terminal
}

func (authCtx Authorization) TerminalID() TerminalID {
	return authCtx.Session.terminal.id
}

// TerminalIDPtr returns a pointer to a new copy of terminal ID. The
// returned value is non-nil when the terminal ID is valid.
func (authCtx Authorization) TerminalIDPtr() *TerminalID {
	return authCtx.Session.terminal.IDPtr()
}

func (authCtx Authorization) ClientID() ApplicationID {
	return authCtx.Session.terminal.application.ID()
}

// RawToken returns the token where this instance of Authorization
// was parsed from.
func (authCtx Authorization) RawToken() string {
	return authCtx.rawToken
}

const (
	// AccessTokenTTLDefault is the active duration for an access token.
	//
	// We might want to make this configurable.
	AccessTokenTTLDefault = 20 * time.Minute
	// AccessTokenTTLDefaultInSeconds is a shortcut to get AccessTokenTTLDefault in seconds.
	AccessTokenTTLDefaultInSeconds = int64(AccessTokenTTLDefault / time.Second)
)

type AccessTokenClaims struct {
	jwt.Claims

	AuthorizedParty string `json:"azp,omitempty"`
	SubType         string `json:"sub_type,omitempty"`
	TerminalID      string `json:"terminal_id,omitempty"`
}

//TODO: unused. remove this.
func (claims AccessTokenClaims) Valid() error {
	if claims.ID != "" {
		return nil
	}
	return errors.EntMsg("jti", "empty")
}

// RefreshTokenTTLDefault is the active duration for a refresh token.
//
// We might want to make this configurable.
const RefreshTokenTTLDefault = 30 * 24 * time.Hour

type RefreshTokenClaims struct {
	ExpiresAt      int64  `json:"exp,omitempty"`
	NotBefore      int64  `json:"nbf,omitempty"`
	TerminalID     string `json:"terminal_id,omitempty"`
	TerminalSecret string `json:"terminal_secret,omitempty"`
}

// Valid is provided as required for claims. Do not use this method.
func (claims RefreshTokenClaims) Valid() error {
	return nil
}
