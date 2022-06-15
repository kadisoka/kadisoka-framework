package iam

import (
	"time"

	"github.com/alloyzeus/go-azfl/azcore"
	"github.com/alloyzeus/go-azfl/errors"
	dataerrs "github.com/alloyzeus/go-azfl/errors/data"
	"github.com/square/go-jose/v3/jwt"
)

// Used in API call metadata: HTTP header and gRPC call metadata
const (
	AuthorizationMetadataKey    = "Authorization"
	AuthorizationMetadataKeyAlt = "authorization"
)

var (
	ErrReqFieldAuthorizationMalformed = ReqFieldErr(AuthorizationMetadataKey, dataerrs.ErrMalformed)

	ErrReqFieldAuthorizationTypeUnsupported = ReqFieldErr(AuthorizationMetadataKey, dataerrs.ErrTypeUnsupported)

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

	SessionID SessionID

	// Scope, expiry time

	rawToken string
}

var _ azcore.Session[
	SessionIDNum, SessionID,
	TerminalIDNum, TerminalID,
	UserIDNum, UserID,
	AuthorizationSubject,
] = Authorization{}

func (authz Authorization) ParentSessionID() SessionID {
	if authz.AssumingAuthorization != nil {
		return authz.AssumingAuthorization.SessionID
	}
	return SessionIDZero()
}

func (authz Authorization) ID() SessionID { return authz.SessionID }

func (authz Authorization) Subject() AuthorizationSubject {
	return NewAuthorizationSubject(
		authz.SessionID.terminal, authz.SessionID.terminal.user)
}

// newEmptyAuthorization creates a new instance of Authorization without
// any data.
func newEmptyAuthorization() *Authorization {
	return &Authorization{}
}

func (authz Authorization) IsStaticallyValid() bool {
	return authz.SessionID.IsStaticallyValid()
}

func (authz Authorization) IsNotStaticallyValid() bool {
	return !authz.IsStaticallyValid()
}

// IsTerminal returns true if the authorized terminal is the same as terminalID.
func (authz Authorization) IsTerminal(terminalID TerminalID) bool {
	ctxTerm := authz.SessionID.terminal
	return ctxTerm.IsStaticallyValid() && ctxTerm.EqualsTerminalID(terminalID)
}

// IsUser checks if this authorization is represeting a particular user.
func (authz Authorization) IsUser(userID UserID) bool {
	return authz.ClientApplicationIDNum().IsUserAgent() &&
		authz.SessionID.terminal.user.EqualsUserID(userID)
}

// IsUserSubject is used to determine if this authorization represents a user.
func (authz Authorization) IsUserSubject() bool {
	if authz.ClientApplicationIDNum().IsUserAgent() &&
		authz.SessionID.terminal.user.IsStaticallyValid() {
		return true
	}
	return false
}

func (authz Authorization) IsServiceClientContext() bool {
	if authz.ClientApplicationIDNum().IsService() &&
		authz.SessionID.terminal.user.IsNotStaticallyValid() {
		return true
	}
	return false
}

func (authz Authorization) UserID() UserID {
	return authz.SessionID.terminal.user
}

// UserIDPtr returns a pointer to a new copy of user ref-key. The
// returned value is non-nil when the user ref-key is valid.
func (authz Authorization) UserIDPtr() *UserID {
	return authz.SessionID.terminal.UserPtr()
}

func (authz Authorization) UserIDNum() UserIDNum {
	return authz.SessionID.terminal.user.IDNum()
}

// UserIDNumPtr returns a pointer to a new copy of user id-num. The
// returned value is non-nil when the user id-num is valid.
func (authz Authorization) UserIDNumPtr() *UserIDNum {
	return authz.SessionID.terminal.user.IDNumPtr()
}

func (authz Authorization) TerminalID() TerminalID {
	return authz.SessionID.terminal
}

func (authz Authorization) TerminalIDNum() TerminalIDNum {
	return authz.SessionID.terminal.idNum
}

// TerminalIDNumPtr returns a pointer to a new copy of terminal id-num. The
// returned value is non-nil when the terminal id-num is valid.
func (authz Authorization) TerminalIDNumPtr() *TerminalIDNum {
	return authz.SessionID.terminal.IDNumPtr()
}

func (authz Authorization) ClientApplicationIDNum() ApplicationIDNum {
	return authz.SessionID.terminal.application.IDNum()
}

// RawToken returns the token where this instance of Authorization
// was parsed from.
func (authz Authorization) RawToken() string {
	return authz.rawToken
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
