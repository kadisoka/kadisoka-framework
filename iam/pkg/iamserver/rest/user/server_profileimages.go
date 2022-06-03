package user

import (
	"net/http"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/emicklei/go-restful"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
)

type userProfileImagePutResponse struct {
	URL string `json:"url"`
}

const multipartFormMaxMemory = 20 * 1024 * 1024

func (restSrv *Server) putUserProfileImage(req *restful.Request, resp *restful.Response) {
	reqCtx, err := restSrv.RESTOpInputContext(req.Request)
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

	if err := req.Request.ParseMultipartForm(multipartFormMaxMemory); err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Form data parsing")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	uploadedFile, _, err := req.Request.FormFile("body")
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request file retrieval")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}
	defer uploadedFile.Close()

	imageURL, err := restSrv.serverCore.
		SetUserProfileImageByFile(reqCtx, ctxAuth.UserRef(), uploadedFile)
	if err != nil {
		if errors.IsCallError(err) {
			//TODO: translate the error
			logCtx(reqCtx).
				Warn().Err(err).
				Msg("User profile image update")
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msg("Unable to update user profile image")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	rest.RespondTo(resp).Success(
		&userProfileImagePutResponse{
			URL: imageURL,
		})
}
