package rest

import (
	"net/http"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api"
)

type RequestContext interface {
	api.CallContext
	HTTPRequest() *http.Request
}
