package iam

import (
	"context"
	"time"

	accesserrs "github.com/alloyzeus/go-azfl/azfl/errors/access"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api"
)

var (
	ErrAuthorizationRequired = accesserrs.Msg("authorization context required")
	ErrAuthorizationInvalid  = accesserrs.Msg("authorization invalid")

	ErrUserContextRequired          = accesserrs.Msg("user context required")
	ErrServiceClientContextRequired = accesserrs.Msg("service client context required")

	ErrContextUserNotAllowedToPerformActionOnResource = accesserrs.Msg("context user is not allowed perform action on the target resource")
	ErrContextUserNotAllowedToAccessToOthersResource  = accesserrs.Msg("context user is not allowed to access to other's resource")
)

func NewEmptyCallContext(ctx context.Context) CallContext {
	return &callContext{
		Context:            ctx,
		authorization:      newEmptyAuthorization(),
		requestReceiveTime: time.Now().UTC(),
	}
}

// CallContext provides call-scoped information.
type CallContext interface {
	api.CallContext
	Authorization() Authorization
	IsUserContext() bool
}

func newCallContext(
	ctx context.Context,
	authCtx *Authorization,
	remoteAddress string,
	remoteEnvString string,
	requestID *api.RequestID,
) CallContext {
	if authCtx == nil {
		panic("authCtx must not be nil")
	}
	return &callContext{ctx, authCtx, remoteAddress, remoteEnvString,
		requestID, time.Now().UTC()}
}

var _ CallContext = &callContext{}

type callContext struct {
	context.Context
	authorization      *Authorization
	remoteAddress      string
	remoteEnvString    string
	requestID          *api.RequestID
	requestReceiveTime time.Time
}

func (ctx callContext) Authorization() Authorization {
	if ctx.authorization == nil {
		authCtx := newEmptyAuthorization()
		return *authCtx
	}
	return *ctx.authorization
}

func (ctx *callContext) IsUserContext() bool {
	return ctx != nil && ctx.authorization != nil &&
		ctx.authorization.IsUserContext()
}

func (ctx *callContext) MethodName() string { return "" }

func (ctx *callContext) RequestID() *api.RequestID { return ctx.requestID }

func (ctx *callContext) RemoteAddress() string { return ctx.remoteAddress }

func (ctx *callContext) RemoteEnvironmentString() string { return ctx.remoteEnvString }

func (ctx *callContext) RequestReceiveTime() time.Time { return ctx.requestReceiveTime }
