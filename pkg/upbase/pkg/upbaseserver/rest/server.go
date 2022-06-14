package rest

import (
	"net/http"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam/rest/logging"
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

func (restSrv *Server) RESTCallInputContext(req *http.Request) (*upbase.RESTCallInputContext, error) {
	return restSrv.serverCore.RESTCallInputContext(req)
}
