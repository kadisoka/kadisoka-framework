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

	ErrOperationContextMissing      = accesserrs.Msg("operation context is missing")
	ErrUserContextRequired          = accesserrs.Msg("user context required")
	ErrServiceClientContextRequired = accesserrs.Msg("service client context required")

	ErrOperationNotAllowed = accesserrs.Msg("actor is not allowed perform action on the target resource")
	ErrAccessNotAllowed    = accesserrs.Msg("actor is not allowed to access target resource")
)

func NewEmptyOpInputContext(ctx context.Context) OpInputContext {
	return &callContext{
		Context:       ctx,
		authorization: newEmptyAuthorization(),
		requestInfo: api.OpInputMetadata{
			ReceiveTime: time.Now().UTC(),
		},
	}
}

// OpInputContext provides call-scoped information.
type OpInputContext interface {
	api.OpInputContext[
		SessionIDNum, SessionRefKey, TerminalIDNum, TerminalRefKey,
		UserIDNum, UserRefKey, Actor, Authorization]

	Authorization() Authorization
	IsUserContext() bool
}

func newOpInputContext(
	ctx context.Context,
	ctxAuth *Authorization,
	originInfo azcore.ServiceMethodCallOriginInfo,
	requestID *api.OpID,
) OpInputContext {
	if ctxAuth == nil {
		panic("ctxAuth must not be nil")
	}
	return &callContext{ctx, ctxAuth, api.OpInputMetadata{
		ID:          requestID,
		ReceiveTime: time.Now().UTC(),
	}, originInfo,
	}
}

var _ OpInputContext = &callContext{}

type callContext struct {
	context.Context
	authorization *Authorization
	requestInfo   api.OpInputMetadata
	originInfo    azcore.ServiceMethodCallOriginInfo
}

func (callContext) AZContext()                       {}
func (callContext) AZServiceContext()                {}
func (callContext) AZServiceMethodContext()          {}
func (callContext) AZServiceMethodCallContext()      {}
func (callContext) AZServiceMethodCallInputContext() {}
func (ctx callContext) ServiceMethodCallOriginInfo() azcore.ServiceMethodCallOriginInfo {
	return ctx.originInfo
}
func (ctx callContext) Session() Authorization {
	return *ctx.authorization
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

func (ctx *callContext) OpInputMetadata() api.OpInputMetadata { return ctx.requestInfo }

type OpOutputContext struct {
	Err     error
	Mutated bool
}
