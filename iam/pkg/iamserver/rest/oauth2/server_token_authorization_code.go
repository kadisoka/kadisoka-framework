//

package oauth2

import (
	"strings"
	"time"

	"github.com/alloyzeus/go-azcore/azcore/errors"
	"github.com/emicklei/go-restful"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/oauth2"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (restSrv *Server) handleTokenRequestByAuthorizationCodeGrant(
	req *restful.Request, resp *restful.Response,
) {
	reqClient, err := restSrv.serverCore.
		RequestClient(req.Request)
	if reqClient == nil {
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

	var termRef iam.TerminalRefKey
	if strings.HasPrefix(authCode, "otp:") {
		// Only for non-confidential user-agents
		if appRef := reqClient.ID; !appRef.ID().IsUserAgentAuthorizationPublic() {
			logReq(req.Request).
				Warn().Msgf("Client %v is not allowed to use grant type 'authorization_code'", reqClient.ID)
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorUnauthorizedClient)
			return
		}

		parts := strings.Split(authCode, ":")
		if len(parts) != 3 {
			logReq(req.Request).
				Warn().Msgf("Authorization code contains invalid number of parts (%v)", len(parts))
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		termRef, err = iam.TerminalRefKeyFromAZERText(parts[1])
		if err != nil || termRef.IsNotValid() {
			logReq(req.Request).
				Warn().Err(err).Msg("Auth code malformed")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		authCode = parts[2]
	} else {
		// Only for confidential user-agents
		if appRef := reqClient.ID; !appRef.ID().IsUserAgentAuthorizationConfidential() {
			logReq(req.Request).
				Warn().Msgf("Client %v is not allowed to use grant type 'authorization_code'", reqClient.ID)
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorUnauthorizedClient)
			return
		}

		termRef, err = iam.TerminalRefKeyFromAZERText(authCode)
		if err != nil || termRef.IsNotValid() {
			logReq(req.Request).
				Warn().Err(err).Msg("Auth code malformed")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		authCode = ""
	}

	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil && err != iam.ErrReqFieldAuthorizationTypeUnsupported {
		logCtx(reqCtx).
			Warn().Err(err).Msg("Request context")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}
	authCtx := reqCtx.Authorization()
	if authCtx.IsValid() {
		logCtx(reqCtx).
			Warn().Msg("Authorization context must not be valid")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	redirectURI := req.Request.FormValue("redirect_uri")
	if redirectURI != "" && reqClient.HasOAuth2RedirectURI(redirectURI) {
		logCtx(reqCtx).
			Warn().Msgf("Invalid redirect_uri: %s (wants %s)", redirectURI, reqClient.OAuth2RedirectURI)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidRequest)
		return
	}

	clientAZERText := req.Request.FormValue("client_id")
	if clientAZERText != "" && clientAZERText != reqClient.ID.AZERText() {
		logCtx(reqCtx).
			Warn().Msgf("Invalid client_id: %s (wants %s)", clientAZERText, reqClient.ID)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidClient)
		return
	}

	terminalSecret, userRef, err := restSrv.serverCore.
		ConfirmTerminalRegistrationVerification(reqCtx, termRef, authCode)
	if err != nil {
		switch err {
		case iam.ErrTerminalVerificationCodeExpired:
			logCtx(reqCtx).
				Warn().Err(err).Msg("ConfirmTerminalAuthorization")
			// Status code 410 (gone) might be more approriate but the standard
			// says that we should use 400 for expired grant.
			oauth2.RespondTo(resp).Error(oauth2.ErrorResponse{
				Error:            oauth2.ErrorInvalidGrant,
				ErrorDescription: "expired"})
			return
		case iam.ErrAuthorizationCodeAlreadyClaimed,
			iam.ErrTerminalVerificationCodeMismatch:
			logCtx(reqCtx).
				Warn().Err(err).Msg("ConfirmTerminalAuthorization")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).Msg("ConfirmTerminalAuthorization")
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidRequest)
			return
		}
		logCtx(reqCtx).
			Err(err).Msgf("ConfirmTerminalAuthorization")
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
		logCtx(reqCtx).
			Error().Msgf("GenerateRefreshTokenJWT: %v", err)
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
			UserID:         userRef.AZERText(),
			TerminalSecret: terminalSecret,
		})
}
