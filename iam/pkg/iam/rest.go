package iam

import (
	"net/http"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
)

// RESTServiceClient is the interface specialized for REST.
type RESTServiceClient interface {
	// AuthorizedOutgoingHTTPRequestHeader returns a new instance of http.Header
	// with authorization information set. If baseHeader is proivded, this method
	// will merge it into the returned value.
	AuthorizedOutgoingHTTPRequestHeader(
		baseHeader http.Header,
	) http.Header
}

type RESTCallInputContext struct {
	CallInputContext

	Request *http.Request
}

var _ rest.CallInputContext[
	SessionIDNum, SessionRefKey, TerminalIDNum, TerminalRefKey,
	UserIDNum, UserRefKey, Actor, Authorization, api.IdempotencyKey,
] = &RESTCallInputContext{}

func (reqCtx *RESTCallInputContext) HTTPRequest() *http.Request {
	if reqCtx != nil {
		return reqCtx.Request
	}
	return nil
}

func (reqCtx *RESTCallInputContext) MethodName() string {
	if reqCtx == nil || reqCtx.Request == nil {
		return ""
	}
	return reqCtx.Request.Method
}

func (reqCtx *RESTCallInputContext) ResourceID() string {
	if reqCtx == nil || reqCtx.Request == nil {
		return ""
	}
	return reqCtx.Request.URL.String()
}
