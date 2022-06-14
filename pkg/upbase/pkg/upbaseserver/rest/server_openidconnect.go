package rest

import (
	"net/http"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/emicklei/go-restful/v3"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/rest"
)

//TODO: the details would be depends on the type of the client:
// if it's internal, it could get all the details. Otherwise, it will
// be depended on the requested scope and user's privacy settings.
//TODO: check error responses spec
func (restSrv *Server) getUserOpenIDConnectUserInfo(
	req *restful.Request, resp *restful.Response,
) {
	reqCtx, err := restSrv.RESTCallInputContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request context")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	ctxAuth := reqCtx.Authorization()

	userInfo, err := restSrv.serverCore.
		GetUserOpenIDConnectStandardClaims(reqCtx, ctxAuth.UserID())
	if err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msg("GetUserOpenIDConnectStandardClaims")
			rest.RespondTo(resp).
				EmptyError(http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msg("GetUserOpenIDConnectStandardClaims")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	restSrv.eTagResponder.RespondGetJSON(req, resp, &userInfo)
}
