package rest

import (
	"net/http"

	"github.com/kadisoka/foundation/pkg/api"
)

type RequestContext interface {
	api.CallContext
	HTTPRequest() *http.Request
}
