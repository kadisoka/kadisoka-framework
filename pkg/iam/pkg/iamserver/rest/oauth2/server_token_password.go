//

package oauth2

import (
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/emicklei/go-restful/v3"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/oauth2"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
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

	if reqApp == nil {
		logReq(req.Request).
			Warn().Str("applicationID", reqApp.ID.AZIDText()).
			Msg("Application authentication is required")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorUnauthorizedClient)
		return
	}

	reqCtx, err := restSrv.RESTCallInputContext(req.Request)
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

	//TODO: move this to core
	// Username with scheme. The format is '<scheme>:<scheme-specific-identifier>'
	if names := strings.SplitN(username, ":", 2); len(names) == 2 {
		switch names[0] {
		case "terminal":
			restSrv.handleTokenRequestByPasswordGrantWithTerminalCreds(
				reqCtx, resp, reqApp, names[1], password)
			return
		}
	}

	termID, termSecret, userID, err := restSrv.serverCore.
		AuthorizeTerminalByUserIdentifierAndPassword(reqCtx, reqApp, "", username, password)
	if err != nil {
		if _, ok := err.(errors.CallError); ok {
			logReq(req.Request).
				Warn().Err(err).
				Msg("AuthorizeTerminalByUserIdentifierAndPassword")
			//TODO: be more accurate about the error
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		logReq(req.Request).
			Error().Err(err).
			Msg("AuthorizeTerminalByUserIdentifierAndPassword")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	if userID.IsNotStaticallyValid() {
		logReq(req.Request).
			Warn().Str("username", username).
			Msg("Authentication failed")
		return
	}

	accessToken, refreshToken, err := restSrv.serverCore.
		GenerateTokenSetJWT(reqCtx, termID, ctxAuth.UserID(), termSecret)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).
			Msg("GenerateTokenSetJWT")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	oauth2.RespondTo(resp).TokenCustom(
		&iam.OAuth2TokenResponse{
			TokenResponse: oauth2.TokenResponse{
				AccessToken:  accessToken,
				TokenType:    oauth2.TokenTypeBearer,
				ExpiresIn:    iam.AccessTokenTTLDefaultInSeconds,
				RefreshToken: refreshToken,
			},
			UserID:         userID.AZIDText(),
			TerminalID:     termID.AZIDText(),
			TerminalSecret: termSecret,
		})
}

func (restSrv *Server) handleTokenRequestByPasswordGrantWithTerminalCreds(
	reqCtx *iam.RESTCallInputContext,
	resp *restful.Response,
	reqApp *iam.Application,
	terminalIDStr string,
	terminalSecret string,
) {
	termID, err := iam.TerminalIDFromAZIDText(terminalIDStr)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("terminalIDStr", terminalIDStr).
			Msg("Unable to parse username as TerminalID")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	if termID.IsNotStaticallyValid() {
		logCtx(reqCtx).
			Warn().Str("terminalIDStr", terminalIDStr).Str("terminalID", termID.AZIDText()).
			Msg("Terminal ref is invalid")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	appID := termID.Application()
	if !appID.IDNum().IsService() {
		logCtx(reqCtx).
			Warn().Str("terminalID", termID.AZIDText()).
			Msg("Application is not allowed to use grant type")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorUnauthorizedClient)
		return
	}

	if !appID.EqualsApplicationID(reqApp.ID) {
		logCtx(reqCtx).
			Warn().Str("terminalID", termID.AZIDText()).Str("applicationRef", appID.AZIDText()).
			Msg("Terminal credentials are that of other application")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	authOK, userID, err := restSrv.serverCore.
		AuthenticateTerminal(termID, terminalSecret)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).Str("terminalID", termID.AZIDText()).
			Msg("AuthenticateTerminal")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}
	if !authOK {
		logCtx(reqCtx).
			Warn().Str("terminalID", termID.AZIDText()).
			Msg("Terminal authentication failed")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	if userID.IsStaticallyValid() {
		logCtx(reqCtx).
			Warn().Str("terminalID", termID.AZIDText()).Str("userID", userID.AZIDText()).
			Msg("Terminal must not be associated to any user")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	accessToken, refreshToken, err := restSrv.serverCore.
		GenerateTokenSetJWT(reqCtx, termID, userID, terminalSecret)
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
			UserID: userID.AZIDText(),
		})
}
