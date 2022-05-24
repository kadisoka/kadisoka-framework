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
	return ctxAuth.Session.IsStaticallyValid()
}

func (ctxAuth Authorization) IsNotValid() bool {
	return !ctxAuth.IsValid()
}

// IsTerminal returns true if the authorized terminal is the same as termRef.
func (ctxAuth Authorization) IsTerminal(termRef TerminalRefKey) bool {
	ctxTerm := ctxAuth.Session.terminal
	return ctxTerm.IsStaticallyValid() && ctxTerm.EqualsTerminalRefKey(termRef)
}

func (ctxAuth Authorization) Actor() Actor {
	return Actor{
		UserRef:     ctxAuth.Session.terminal.user,
		TerminalRef: ctxAuth.Session.terminal,
	}
}

// IsUser checks if this authorization is represeting a particular user.
func (ctxAuth Authorization) IsUser(userRef UserRefKey) bool {
	return ctxAuth.ClientApplicationIDNum().IsUserAgent() &&
		ctxAuth.Session.terminal.user.EqualsUserRefKey(userRef)
}

// IsUserContext is used to determine if this context represents a user.
func (ctxAuth Authorization) IsUserContext() bool {
	if ctxAuth.ClientApplicationIDNum().IsUserAgent() &&
		ctxAuth.Session.terminal.user.IsStaticallyValid() {
		return true
	}
	return false
}

func (ctxAuth Authorization) IsServiceClientContext() bool {
	if ctxAuth.ClientApplicationIDNum().IsService() &&
		ctxAuth.Session.terminal.user.IsNotStaticallyValid() {
		return true
	}
	return false
}

func (ctxAuth Authorization) UserRef() UserRefKey {
	return ctxAuth.Session.terminal.user
}

// UserRefKeyPtr returns a pointer to a new copy of user ref-key. The
// returned value is non-nil when the user ref-key is valid.
func (ctxAuth Authorization) UserRefKeyPtr() *UserRefKey {
	return ctxAuth.Session.terminal.UserPtr()
}

func (ctxAuth Authorization) UserIDNum() UserIDNum {
	return ctxAuth.Session.terminal.user.IDNum()
}

// UserIDNumPtr returns a pointer to a new copy of user id-num. The
// returned value is non-nil when the user id-num is valid.
func (ctxAuth Authorization) UserIDNumPtr() *UserIDNum {
	return ctxAuth.Session.terminal.user.IDNumPtr()
}

func (ctxAuth Authorization) TerminalRef() TerminalRefKey {
	return ctxAuth.Session.terminal
}

func (ctxAuth Authorization) TerminalIDNum() TerminalIDNum {
	return ctxAuth.Session.terminal.idNum
}

// TerminalIDNumPtr returns a pointer to a new copy of terminal id-num. The
// returned value is non-nil when the terminal id-num is valid.
func (ctxAuth Authorization) TerminalIDNumPtr() *TerminalIDNum {
	return ctxAuth.Session.terminal.IDNumPtr()
}

func (ctxAuth Authorization) ClientApplicationIDNum() ApplicationIDNum {
	return ctxAuth.Session.terminal.application.IDNum()
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
