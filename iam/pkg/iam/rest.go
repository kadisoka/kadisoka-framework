package iam

import (
	"net/http"

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

type RESTOpInputContext struct {
	OpInputContext
	Request *http.Request
}

var _ rest.OpInputContext = &RESTOpInputContext{}

func (reqCtx *RESTOpInputContext) HTTPRequest() *http.Request {
	if reqCtx != nil {
		return reqCtx.Request
	}
	return nil
}

func (reqCtx *RESTOpInputContext) MethodName() string {
	if reqCtx == nil || reqCtx.Request == nil {
		return ""
	}
	req := reqCtx.Request
	var urlStr string
	if req.URL != nil {
		urlStr = req.URL.String()
	}
	return req.Method + " " + urlStr
}
