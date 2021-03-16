package iam

import (
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	dataerrs "github.com/alloyzeus/go-azfl/azfl/errors/data"
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
// token (access token) provided along the request/call.
//TODO: include the application ref if it's using client authentication.
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

func (ctxAuth Authorization) IsValid() bool {
	return ctxAuth.Session.IsValid()
}

func (ctxAuth Authorization) IsNotValid() bool {
	return !ctxAuth.IsValid()
}

// IsTerminal returns true if the authorized terminal is the same as termRef.
func (ctxAuth Authorization) IsTerminal(termRef TerminalRefKey) bool {
	ctxTerm := ctxAuth.Session.terminal
	return ctxTerm.IsValid() && ctxTerm.EqualsTerminalRefKey(termRef)
}

func (ctxAuth Authorization) Actor() Actor {
	return Actor{
		UserRef:     ctxAuth.Session.terminal.user,
		TerminalRef: ctxAuth.Session.terminal,
	}
}

// IsUser checks if this authorization is represeting a particular user.
func (ctxAuth Authorization) IsUser(userRef UserRefKey) bool {
	return ctxAuth.ClientID().IsUserAgent() && ctxAuth.Session.terminal.user.EqualsUserRefKey(userRef)
}

// IsUserContext is used to determine if this context represents a user.
func (ctxAuth Authorization) IsUserContext() bool {
	if ctxAuth.ClientID().IsUserAgent() && ctxAuth.Session.terminal.user.IsValid() {
		return true
	}
	return false
}

func (ctxAuth Authorization) IsServiceClientContext() bool {
	if ctxAuth.ClientID().IsService() && ctxAuth.Session.terminal.user.IsNotValid() {
		return true
	}
	return false
}

func (ctxAuth Authorization) UserRef() UserRefKey {
	return ctxAuth.Session.terminal.user
}

// UserRefKeyPtr returns a pointer to a new copy of user ID. The
// returned value is non-nil when the user ref-key is valid.
func (ctxAuth Authorization) UserRefKeyPtr() *UserRefKey {
	if ctxAuth.Session.terminal.user.IsValid() {
		return &ctxAuth.Session.terminal.user
	}
	return nil
}

func (ctxAuth Authorization) UserID() UserID {
	return ctxAuth.Session.terminal.user.ID()
}

// UserIDPtr returns a pointer to a new copy of user ID. The
// returned value is non-nil when the user ref-key is valid.
func (ctxAuth Authorization) UserIDPtr() *UserID {
	return ctxAuth.Session.terminal.user.IDPtr()
}

func (ctxAuth Authorization) TerminalRef() TerminalRefKey {
	return ctxAuth.Session.terminal
}

func (ctxAuth Authorization) TerminalID() TerminalID {
	return ctxAuth.Session.terminal.id
}

// TerminalIDPtr returns a pointer to a new copy of terminal ID. The
// returned value is non-nil when the terminal ID is valid.
func (ctxAuth Authorization) TerminalIDPtr() *TerminalID {
	return ctxAuth.Session.terminal.IDPtr()
}

func (ctxAuth Authorization) ClientID() ApplicationID {
	return ctxAuth.Session.terminal.application.ID()
}

// RawToken returns the token where this instance of Authorization
// was parsed from.
func (ctxAuth Authorization) RawToken() string {
	return ctxAuth.rawToken
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
