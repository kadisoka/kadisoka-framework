//

package oauth2

import (
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/emicklei/go-restful/v3"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/oauth2"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

func (restSrv *Server) handleTokenRequestByAuthorizationCodeGrant(
	req *restful.Request, resp *restful.Response,
) {
	reqApp, err := restSrv.serverCore.
		RequestApplication(req.Request)
	if reqApp == nil {
		if err != nil {
			logReq(req.Request).
				Warn().Err(err).Msg("Client authentication")
		} else {
			logReq(req.Request).
				Warn().Msg("No authorized client")
		}
		// RFC 6749 § 5.2
		oauth2.RespondTo(resp).ErrInvalidClientBasicAuthorization(
			restSrv.serverCore.RealmName(), "")
		return
	}

	authCode := req.Request.FormValue("code")
	if authCode == "" {
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	var termID iam.TerminalID
	if strings.HasPrefix(authCode, "otp:") {
		// Only for non-confidential user-agents
		if appID := reqApp.ID; !appID.IDNum().IsUserAgentAuthorizationPublic() {
			logReq(req.Request).
				Warn().Str("client_id", reqApp.ID.AZIDText()).
				Msg("Client is not allowed to use grant type 'authorization_code' with OTP")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorUnauthorizedClient)
			return
		}

		parts := strings.Split(authCode, ":")
		if len(parts) != 3 {
			logReq(req.Request).
				Warn().Str("code", authCode).
				Msg("Code contains invalid number of parts")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		termIDStr := parts[1]
		termID, err = iam.TerminalIDFromAZIDText(termIDStr)
		if err != nil || termID.IsNotStaticallyValid() {
			logReq(req.Request).
				Warn().Err(err).Str("code", authCode).
				Msg("Code malformed")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		authCode = parts[2]
	} else {
		// Only for confidential user-agents
		if appID := reqApp.ID; !appID.IDNum().IsUserAgentAuthorizationConfidential() {
			logReq(req.Request).
				Warn().Str("client_id", reqApp.ID.AZIDText()).
				Msg("Client is not allowed to use grant type 'authorization_code'")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorUnauthorizedClient)
			return
		}

		termID, err = iam.TerminalIDFromAZIDText(authCode)
		if err != nil || termID.IsNotStaticallyValid() {
			logReq(req.Request).
				Warn().Err(err).Str("code", authCode).
				Msg("Code malformed")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		authCode = ""
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

	redirectURI := req.Request.FormValue("redirect_uri")
	if redirectURI != "" && reqApp.Attributes.HasOAuth2RedirectURI(redirectURI) {
		logCtx(reqCtx).
			Warn().Msgf("Invalid redirect_uri: %s (wants %s)",
			redirectURI, reqApp.Attributes.OAuth2RedirectURI)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidRequest)
		return
	}

	clientAZIDText := req.Request.FormValue("client_id")
	if clientAZIDText != "" && clientAZIDText != reqApp.ID.AZIDText() {
		logCtx(reqCtx).
			Warn().Msgf("Invalid client_id: %s (wants %s)", clientAZIDText, reqApp.ID)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidClient)
		return
	}

	terminalSecret, userID, err := restSrv.serverCore.
		ConfirmTerminalAuthorization(reqCtx, termID, authCode)
	if err != nil {
		switch err {
		case iam.ErrTerminalVerificationCodeExpired:
			logCtx(reqCtx).
				Warn().Err(err).
				Msg("ConfirmTerminalAuthorization")
			// Status code 410 (gone) might be more approriate but the standard
			// says that we should use 400 for expired grant.
			oauth2.RespondTo(resp).Error(oauth2.ErrorResponse{
				Error:            oauth2.ErrorInvalidGrant,
				ErrorDescription: "expired"})
			return
		case iam.ErrAuthorizationCodeAlreadyClaimed,
			iam.ErrTerminalVerificationCodeMismatch:
			logCtx(reqCtx).
				Warn().Err(err).
				Msg("ConfirmTerminalAuthorization")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msg("ConfirmTerminalAuthorization")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidRequest)
			return
		}
		logCtx(reqCtx).
			Warn().Err(err).
			Msgf("ConfirmTerminalAuthorization")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	accessToken, refreshToken, err := restSrv.serverCore.
		GenerateTokenSetJWT(reqCtx, termID, userID, terminalSecret)
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
			TerminalSecret: terminalSecret,
		})
}
