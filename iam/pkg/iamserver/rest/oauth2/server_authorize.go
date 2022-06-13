//

package oauth2

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/emicklei/go-restful/v3"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/oauth2"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
)

func (restSrv *Server) getAuthorize(req *restful.Request, resp *restful.Response) {
	//TODO: if authorization context is valid, and the application has been
	// previously authorized for the user, simply redirect back.

	r := req.Request
	w := resp

	inQuery := r.URL.Query()
	val, err := oauth2.AuthorizationRequestFromURLValues(inQuery)
	if err != nil {
		logReq(r).
			Error().Err(err).
			Msg("Unable to decode query")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	//TODO:
	// - note that redirect_uri is allowed to be empty
	// - if redirect_uri is not empty, load client data and compare the
	//   redirect_uri. if they are not equal, that's an error
	// - if provided redirect_uri is empty, use client's data
	// - if we have no valid redirect_uri, show error page

	//TODO: support OOB redirect scheme
	if val.RedirectURI != "" && !strings.HasPrefix(val.RedirectURI, "http") {
		logReq(r).
			Warn().Msg("redirect_uri invalid")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	//TODO: validate inputs
	if val.ClientID == "" {
		if val.RedirectURI == "" {
			logReq(r).
				Warn().Msg("client_id invalid and no redirect_uri")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}
		logReq(r).
			Warn().Msg("client_id missing")
		cbURL := val.RedirectURI + "?" + oauth2.MustQueryString(oauth2.ErrorResponse{
			Error: oauth2.ErrorInvalidRequest,
			State: val.State,
		})
		http.Redirect(w, r, cbURL, http.StatusFound)
		return
	}

	appID, err := iam.ApplicationIDFromAZIDText(val.ClientID)
	if err != nil {
		logReq(r).
			Warn().Err(err).
			Msg("client_id malformed")
		cbURL := val.RedirectURI + "?" + oauth2.MustQueryString(oauth2.ErrorResponse{
			Error: oauth2.ErrorInvalidRequest,
			State: val.State,
		})
		http.Redirect(w, r, cbURL, http.StatusFound)
		return
	}
	if appID.IsNotStaticallyValid() {
		logReq(r).
			Warn().Err(err).
			Msg("client_id is invalid")
		cbURL := val.RedirectURI + "?" + oauth2.MustQueryString(oauth2.ErrorResponse{
			Error: oauth2.ErrorInvalidRequest,
			State: val.State,
		})
		http.Redirect(w, r, cbURL, http.StatusFound)
		return
	}

	app, err := restSrv.serverCore.ApplicationByID(appID)
	if err != nil || app == nil {
		logReq(r).
			Warn().Err(err).
			Msg("client_id does not refer a valid client")
		cbURL := val.RedirectURI + "?" + oauth2.MustQueryString(oauth2.ErrorResponse{
			Error: oauth2.ErrorInvalidRequest,
			State: val.State,
		})
		http.Redirect(w, r, cbURL, http.StatusFound)
		return
	}
	if val.RedirectURI != "" && !app.Attributes.HasOAuth2RedirectURI(val.RedirectURI) {
		logReq(r).
			Warn().Msgf("redirect_uri unrecognized %v", val.RedirectURI)
		cbURL := val.RedirectURI + "?" + oauth2.MustQueryString(oauth2.ErrorResponse{
			Error: oauth2.ErrorInvalidRequest,
			State: val.State,
		})
		http.Redirect(w, r, cbURL, http.StatusFound)
		return
	}

	responseTypeArgVal, _ := req.BodyParameter("response_type")
	responseType := oauth2.ResponseTypeFromString(responseTypeArgVal)
	if responseType != oauth2.ResponseTypeCode {
		logReq(r).
			Warn().Str("query.response_type", responseType.String()).
			Msg("Unsupported")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	//TODO:
	// - check the scopes
	// - ensure that the client is allowed to use this flow

	targetURL := restSrv.signInURL + "?" + inQuery.Encode()
	http.Redirect(w, r, targetURL, http.StatusFound)
}

//TODO: some stuff should be moved into core
func (restSrv *Server) postAuthorize(req *restful.Request, resp *restful.Response) {
	reqCtx, err := restSrv.RESTCallInputContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request context")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	ctxAuth := reqCtx.Authorization()
	if !ctxAuth.IsUserSubject() {
		logCtx(reqCtx).
			Warn().Msg("User context required")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	appIDArgVal, _ := req.BodyParameter("client_id")
	appID, err := iam.ApplicationIDFromAZIDText(appIDArgVal)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("form.client_id", appIDArgVal).
			Msg("Malformed")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	responseTypeArgVal, _ := req.BodyParameter("response_type")
	responseType := oauth2.ResponseTypeFromString(responseTypeArgVal)
	if responseType != oauth2.ResponseTypeCode {
		logCtx(reqCtx).
			Warn().Str("form.response_type", responseType.String()).
			Msg("Unsupported")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	app, err := restSrv.serverCore.ApplicationByID(appID)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).Str("client_id", appID.AZIDText()).
			Msg("ApplicationByID")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	if app == nil {
		logCtx(reqCtx).
			Warn().Str("client_id", appID.AZIDText()).
			Msg("Not found")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}
	if !appID.IDNum().IsUserAgentAuthorizationConfidential() {
		logCtx(reqCtx).
			Warn().Str("client_id", appID.AZIDText()).
			Msg("Requires ua-confidential client type")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	redirectURIStr, _ := req.BodyParameter("redirect_uri")
	if redirectURIStr != "" && !app.Attributes.HasOAuth2RedirectURI(redirectURIStr) {
		logCtx(reqCtx).
			Warn().Msgf("Redirect URI mismatch for client %v. Got %v , expecting %v .",
			appID, redirectURIStr, app.Attributes.OAuth2RedirectURI)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	if redirectURIStr == "" {
		redirectURIStr = app.Attributes.OAuth2RedirectURI[0]
	}
	redirectURI, err := url.Parse(redirectURIStr)
	if err != nil {
		panic(err)
	}

	state, _ := req.BodyParameter("state")
	termDisplayName := ""
	var termID iam.TerminalID

	regOutCtx, regOutData := restSrv.serverCore.
		RegisterTerminal(reqCtx,
			iamserver.TerminalRegistrationInputData{
				ApplicationID:    appID,
				UserID:           ctxAuth.UserID(),
				DisplayName:      termDisplayName,
				VerificationType: iam.TerminalVerificationResourceTypeOAuthAuthorizationCode,
				VerificationID:   0,
			})
	if err := regOutCtx.Err; err != nil {
		panic(err)
	}

	termID = regOutData.TerminalID

	redirectURI.RawQuery = oauth2.MustQueryString(oauth2.AuthorizationResponse{
		Code:  termID.AZIDText(),
		State: state,
	})

	rest.RespondTo(resp).Success(
		&iam.OAuth2AuthorizePostResponse{
			RedirectURI: redirectURI.String(),
		})
}
