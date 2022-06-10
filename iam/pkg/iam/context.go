package iam

import (
	"context"
	"time"

	"github.com/alloyzeus/go-azfl/azcore"
	accesserrs "github.com/alloyzeus/go-azfl/errors/access"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api"
)

var (
	ErrAuthorizationRequired = accesserrs.Msg("authorization context required")
	ErrAuthorizationInvalid  = accesserrs.Msg("authorization invalid")

	ErrCallInputContextMissing      = accesserrs.Msg("call input context is missing")
	ErrUserContextRequired          = accesserrs.Msg("user context required")
	ErrServiceClientContextRequired = accesserrs.Msg("service client context required")

	ErrOperationNotAllowed = accesserrs.Msg("actor is not allowed perform action on the target resource")
	ErrAccessNotAllowed    = accesserrs.Msg("actor is not allowed to access target resource")
)

func NewEmptyCallInputContext(ctx context.Context) CallInputContext {
	return &callContext{
		Context:       ctx,
		authorization: newEmptyAuthorization(),
		requestInfo: api.CallInputMetadata{
			ReceiveTime: time.Now().UTC(),
		},
	}
}

// CallInputContext provides call-scoped information.
type CallInputContext interface {
	api.CallInputContext[
		SessionIDNum, SessionID, TerminalIDNum, TerminalID,
		UserIDNum, UserID, Actor, Authorization, api.IdempotencyKey]

	Authorization() Authorization
	IsUserContext() bool
}

func newCallInputContext(
	ctx context.Context,
	ctxAuth *Authorization,
	originInfo azcore.ServiceMethodCallOriginInfo,
	idempotencyKey *api.IdempotencyKey,
) CallInputContext {
	if ctxAuth == nil {
		panic("ctxAuth must not be nil")
	}
	return &callContext{ctx, ctxAuth, api.CallInputMetadata{
		IdempotencyKey: idempotencyKey,
		ReceiveTime:    time.Now().UTC(),
	}, originInfo,
	}
}

var _ CallInputContext = &callContext{}

type callContext struct {
	context.Context
	authorization *Authorization
	requestInfo   api.CallInputMetadata
	originInfo    azcore.ServiceMethodCallOriginInfo
}

func (callContext) AZContext()                       {}
func (callContext) AZServiceContext()                {}
func (callContext) AZServiceMethodContext()          {}
func (callContext) AZServiceMethodCallContext()      {}
func (callContext) AZServiceMethodCallInputContext() {}

func (ctx callContext) OriginInfo() azcore.ServiceMethodCallOriginInfo {
	return ctx.originInfo
}
func (ctx callContext) Session() Authorization {
	if authz := ctx.authorization; authz != nil {
		return *authz
	}
	return *newEmptyAuthorization()
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
		ctx.authorization.IsUserSubject()
}

func (ctx *callContext) MethodName() string { return "" }

func (ctx *callContext) ResourceID() string { return "" }

func (ctx *callContext) IdempotencyKey() api.IdempotencyKey {
	if key := ctx.requestInfo.IdempotencyKey; key != nil {
		return *key
	}
	return api.IdempotencyKeyZero()
}

func (ctx *callContext) CallInputMetadata() api.CallInputMetadata { return ctx.requestInfo }

type CallOutputContext struct {
	Err     error
	Mutated bool
}
