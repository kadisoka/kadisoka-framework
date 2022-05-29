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
	api.OpInputContext
	Authorization() Authorization
	IsUserContext() bool
}

func newOpInputContext(
	ctx context.Context,
	ctxAuth *Authorization,
	originInfo api.OpOriginInfo,
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
	originInfo    api.OpOriginInfo
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

func (ctx *callContext) OpName() string { return "" }

func (ctx *callContext) OpInputMetadata() api.OpInputMetadata { return ctx.requestInfo }

func (ctx *callContext) OpOriginInfo() api.OpOriginInfo { return ctx.originInfo }

type OpOutputContext struct {
	Err     error
	Mutated bool
}
