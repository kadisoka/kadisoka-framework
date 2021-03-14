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
		Context:       ctx,
		authorization: newEmptyAuthorization(),
		requestInfo:   api.CallRequestInfo{ReceiveTime: time.Now().UTC()},
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
	ctxAuth *Authorization,
	originInfo api.CallOriginInfo,
	requestID *api.RequestID,
) CallContext {
	if ctxAuth == nil {
		panic("ctxAuth must not be nil")
	}
	return &callContext{ctx, ctxAuth, api.CallRequestInfo{
		ID:          requestID,
		ReceiveTime: time.Now().UTC(),
	}, originInfo,
	}
}

var _ CallContext = &callContext{}

type callContext struct {
	context.Context
	authorization *Authorization
	requestInfo   api.CallRequestInfo
	originInfo    api.CallOriginInfo
}

func (ctx callContext) Authorization() Authorization {
	if ctx.authorization == nil {
		ctxAuth := newEmptyAuthorization()
		return *ctxAuth
	}
	return *ctx.authorization
}

func (ctx *callContext) IsUserContext() bool {
	return ctx != nil && ctx.authorization != nil &&
		ctx.authorization.IsUserContext()
}

func (ctx *callContext) MethodName() string { return "" }

func (ctx *callContext) RequestInfo() api.CallRequestInfo { return ctx.requestInfo }

func (ctx *callContext) OriginInfo() api.CallOriginInfo { return ctx.originInfo }
