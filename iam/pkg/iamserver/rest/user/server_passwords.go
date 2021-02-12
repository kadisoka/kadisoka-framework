package user

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/kadisoka/foundation/pkg/api/rest"
)

type userPasswordPutRequest struct {
	Password    string `json:"password"`
	OldPassword string `json:"old_password,omitempty"`
}

func (restSrv *Server) putUserPassword(req *restful.Request, resp *restful.Response) {
	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Err(err).Msg("Request context")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	authCtx := reqCtx.Authorization()
	if authCtx.IsNotValid() || !authCtx.IsUserContext() {
		logCtx(reqCtx).
			Warn().Err(err).Msg("Unauthorized")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	var reqBody userPasswordPutRequest
	err = req.ReadEntity(&reqBody)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Msg("Request body parsing")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	matched, err := restSrv.serverCore.
		MatchUserPassword(authCtx.UserID, reqBody.OldPassword)
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
		SetUserPassword(reqCtx, authCtx.UserID, password)
	if err != nil {
		logCtx(reqCtx).
			Err(err).Msg("User password update")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	rest.RespondTo(resp).Success(nil)
}
