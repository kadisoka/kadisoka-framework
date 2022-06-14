package user

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/rest"
)

type userPasswordPutRequest struct {
	Password    string `json:"password"`
	OldPassword string `json:"old_password,omitempty"`
}

//TODO:SEC: harded
func (restSrv *Server) putUserPassword(req *restful.Request, resp *restful.Response) {
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
	if ctxAuth.IsNotStaticallyValid() || !ctxAuth.IsUserSubject() {
		logCtx(reqCtx).
			Warn().Msg("Unauthorized")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	var reqBody userPasswordPutRequest
	err = req.ReadEntity(&reqBody)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Msg("Request entity")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	matched, err := restSrv.serverCore.
		MatchUserPassword(ctxAuth.UserID(), reqBody.OldPassword)
	if err != nil {
		logCtx(reqCtx).
			Err(err).Msg("Passwords matching")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	if !matched {
		logCtx(reqCtx).
			Warn().Msg("Passwords mismatch")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	password := reqBody.Password
	if password == "" {
		logCtx(reqCtx).
			Warn().Msg("Password empty")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	err = restSrv.serverCore.
		SetUserPassword(reqCtx, ctxAuth.UserID(), password)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).
			Msg("SetUserPassword")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	rest.RespondTo(resp).Success(nil)
}
