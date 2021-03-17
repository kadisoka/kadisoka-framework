package rest

import (
	"net/http"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/rest/logging"
	"github.com/kadisoka/kadisoka-framework/upbase/pkg/upbase"
	"github.com/kadisoka/kadisoka-framework/upbase/pkg/upbaseserver"
)

var (
	log    = logging.NewPkgLogger()
	logCtx = log.WithContext
)

type Server struct {
	serverCore    *upbaseserver.RESTServiceServerBase
	eTagResponder *rest.ETagResponder
}

func (restSrv *Server) RESTRequestContext(req *http.Request) (*upbase.RESTRequestContext, error) {
	return restSrv.serverCore.RESTRequestContext(req)
}
