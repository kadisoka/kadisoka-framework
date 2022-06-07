//

package oauth2

import (
	"strings"
	"time"

	"github.com/emicklei/go-restful/v3"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/oauth2"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (restSrv *Server) handleTokenRequestByPasswordGrant(
	req *restful.Request, resp *restful.Response,
) {
	reqApp, err := restSrv.serverCore.
		RequestApplication(req.Request)
	if err != nil {
		logReq(req.Request).
			Warn().Err(err).
			Msg("Client authentication")
		// RFC 6749 § 5.2
		oauth2.RespondTo(resp).ErrInvalidClientBasicAuthorization(
			restSrv.serverCore.RealmName(), "")
		return
	}

	if reqApp != nil && !reqApp.RefKey.IDNum().IsService() {
		logReq(req.Request).
			Warn().Str("applicationID", reqApp.RefKey.AZIDText()).
			Msg("Application is not authorized to use grant type password")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorUnauthorizedClient)
		return
	}

	reqCtx, err := restSrv.RESTOpInputContext(req.Request)
	if err != nil && err != iam.ErrReqFieldAuthorizationTypeUnsupported {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request context")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsStaticallyValid() {
		logCtx(reqCtx).
			Warn().Msg("Authorization context must not be valid")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	username := req.Request.FormValue("username")
	if username == "" {
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}
	password := req.Request.FormValue("password")

	// Username with scheme. The format is '<scheme>:<scheme-specific-identifier>'
	if names := strings.SplitN(username, ":", 2); len(names) == 2 {
		switch names[0] {
		case "terminal":
			restSrv.handleTokenRequestByPasswordGrantWithTerminalCreds(
				reqCtx, resp, reqApp, names[1], password)
			return
		default:
			logReq(req.Request).
				Warn().Str("username", username).
				Msg("Unrecognized identifier scheme")
		}
	}

	logReq(req.Request).
		Warn().Msg("Password grant with no scheme.")
	oauth2.RespondTo(resp).ErrorCode(
		oauth2.ErrorInvalidGrant)
}

func (restSrv *Server) handleTokenRequestByPasswordGrantWithTerminalCreds(
	reqCtx *iam.RESTOpInputContext,
	resp *restful.Response,
	reqApp *iam.Application,
	terminalRefStr string,
	terminalSecret string,
) {
	termRef, err := iam.TerminalRefKeyFromAZIDText(terminalRefStr)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("terminalRefStr", terminalRefStr).
			Msg("Unable to parse username as TerminalRefKey")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	if termRef.IsNotStaticallyValid() {
		logCtx(reqCtx).
			Warn().Str("terminalRefStr", terminalRefStr).Str("terminalRef", termRef.AZIDText()).
			Msg("Terminal ref is invalid")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	appRef := termRef.Application()
	if !appRef.IDNum().IsService() {
		logCtx(reqCtx).
			Warn().Str("terminalRef", termRef.AZIDText()).
			Msg("Application is not allowed to use grant type")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorUnauthorizedClient)
		return
	}

	if !appRef.EqualsApplicationRefKey(reqApp.RefKey) {
		logCtx(reqCtx).
			Warn().Str("terminalRef", termRef.AZIDText()).Str("applicationRef", appRef.AZIDText()).
			Msg("Terminal credentials are that of other application")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	authOK, userRef, err := restSrv.serverCore.
		AuthenticateTerminal(termRef, terminalSecret)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).Str("terminalRef", termRef.AZIDText()).
			Msg("AuthenticateTerminal")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}
	if !authOK {
		logCtx(reqCtx).
			Warn().Str("terminalRef", termRef.AZIDText()).
			Msg("Terminal authentication failed")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	if userRef.IsStaticallyValid() {
		logCtx(reqCtx).
			Warn().Str("terminalRef", termRef.AZIDText()).Str("userRef", userRef.AZIDText()).
			Msg("Terminal must not be associated to any user")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	issueTime := time.Now().UTC()

	accessToken, err := restSrv.serverCore.
		GenerateAccessTokenJWT(reqCtx, termRef, userRef, issueTime)
	if err != nil {
		panic(err)
	}

	refreshToken, err := restSrv.serverCore.
		GenerateRefreshTokenJWT(reqCtx, termRef, terminalSecret, issueTime)
	if err != nil {
		panic(err)
	}

	oauth2.RespondTo(resp).TokenCustom(
		&iam.OAuth2TokenResponse{
			TokenResponse: oauth2.TokenResponse{
				AccessToken:  accessToken,
				TokenType:    oauth2.TokenTypeBearer,
				ExpiresIn:    iam.AccessTokenTTLDefaultInSeconds,
				RefreshToken: refreshToken,
			},
			UserID: userRef.AZIDText(),
		})
}
