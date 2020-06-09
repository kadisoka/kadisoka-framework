package rest

import (
	"net/http"

	"github.com/citadelium/pkg/api"
)

type RequestContext interface {
	api.CallContext
	HTTPRequest() *http.Request
}
