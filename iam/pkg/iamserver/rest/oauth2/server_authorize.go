//

package oauth2

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/emicklei/go-restful"

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

	appRef, err := iam.ApplicationRefKeyFromAZERText(val.ClientID)
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
	if appRef.IsNotValid() {
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
	client, err := restSrv.serverCore.ApplicationByRefKey(appRef)
	if err != nil || client == nil {
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
	if val.RedirectURI != "" && !client.Data.HasOAuth2RedirectURI(val.RedirectURI) {
		logReq(r).
			Warn().Msgf("redirect_uri unrecognized %v", val.RedirectURI)
		cbURL := val.RedirectURI + "?" + oauth2.MustQueryString(oauth2.ErrorResponse{
			Error: oauth2.ErrorInvalidRequest,
			State: val.State,
		})
		http.Redirect(w, r, cbURL, http.StatusFound)
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
	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request context")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	ctxAuth := reqCtx.Authorization()
	if !ctxAuth.IsUserContext() {
		logCtx(reqCtx).
			Warn().Msg("User context required")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	appRefArgVal, _ := req.BodyParameter("client_id")
	appRef, err := iam.ApplicationRefKeyFromAZERText(appRefArgVal)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("form.client_id", appRefArgVal).
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

	client, err := restSrv.serverCore.ApplicationByRefKey(appRef)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).Str("client_id", appRef.AZERText()).
			Msg("ApplicationByRefKey")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	if client == nil {
		logCtx(reqCtx).
			Warn().Str("client_id", appRef.AZERText()).
			Msg("Not found")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}
	if !appRef.ID().IsUserAgentAuthorizationConfidential() {
		logCtx(reqCtx).
			Warn().Str("client_id", appRef.AZERText()).
			Msg("Requires ua-confidential client type")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	redirectURIStr, _ := req.BodyParameter("redirect_uri")
	if redirectURIStr != "" && !client.Data.HasOAuth2RedirectURI(redirectURIStr) {
		logCtx(reqCtx).
			Warn().Msgf("Redirect URI mismatch for client %v. Got %v , expecting %v .",
			appRef, redirectURIStr, client.Data.OAuth2RedirectURI)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	if redirectURIStr == "" {
		redirectURIStr = client.Data.OAuth2RedirectURI[0]
	}
	redirectURI, err := url.Parse(redirectURIStr)
	if err != nil {
		panic(err)
	}

	state, _ := req.BodyParameter("state")
	preferredLanguages := restSrv.parseRequestAcceptLanguage(req, reqCtx)
	termDisplayName := ""
	var termRef iam.TerminalRefKey

	switch responseType {
	case oauth2.ResponseTypeCode:
		termRef, _, err = restSrv.serverCore.
			RegisterTerminal(reqCtx,
				iamserver.TerminalRegistrationInput{
					ApplicationRef:   appRef,
					UserRef:          ctxAuth.UserRef(),
					DisplayName:      termDisplayName,
					AcceptLanguage:   preferredLanguages,
					VerificationType: iam.TerminalVerificationResourceTypeOAuthAuthorizationCode,
					VerificationID:   0,
				})
		if err != nil {
			panic(err)
		}

		redirectURI.RawQuery = oauth2.MustQueryString(oauth2.AuthorizationResponse{
			Code:  termRef.AZERText(),
			State: state,
		})

	case oauth2.ResponseTypeToken:
		termRef, _, err = restSrv.serverCore.
			RegisterTerminal(reqCtx,
				iamserver.TerminalRegistrationInput{
					ApplicationRef:   appRef,
					UserRef:          ctxAuth.UserRef(),
					DisplayName:      termDisplayName,
					AcceptLanguage:   preferredLanguages,
					VerificationType: iam.TerminalVerificationResourceTypeOAuthImplicit,
					VerificationID:   0,
				})
		if err != nil {
			panic(err)
		}

		issueTime := time.Now().UTC()

		tokenString, err := restSrv.serverCore.
			GenerateAccessTokenJWT(reqCtx, termRef, ctxAuth.UserRef(), issueTime)
		if err != nil {
			panic(err)
		}

		redirectURI.Fragment = oauth2.MustQueryString(iam.OAuth2TokenResponse{
			TokenResponse: oauth2.TokenResponse{
				TokenType:   oauth2.TokenTypeBearer,
				AccessToken: tokenString,
				State:       state,
			}})
	}

	rest.RespondTo(resp).Success(
		&iam.OAuth2AuthorizePostResponse{
			RedirectURI: redirectURI.String(),
		})
}
