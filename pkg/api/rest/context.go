package rest

import (
	"net/http"

	"github.com/citadelium/foundation/pkg/api"
)

type RequestContext interface {
	api.CallContext
	HTTPRequest() *http.Request
}
