package rest

import (
	"net/http"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api"
)

type OpInputContext interface {
	api.OpInputContext
	HTTPRequest() *http.Request
}
