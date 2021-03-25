//

package oauth2

import (
	"time"

	"github.com/emicklei/go-restful"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/oauth2"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
)

func (restSrv *Server) handleTokenRequestByClientCredentials(
	req *restful.Request, resp *restful.Response,
) {
	reqApp, err := restSrv.serverCore.
		RequestApplication(req.Request)
	if reqApp == nil {
		if err != nil {
			logReq(req.Request).
				Warn().Err(err).
				Msg("Client authentication")
		} else {
			logReq(req.Request).
				Warn().Msg("No authorized client")
		}
		// RFC 6749 § 5.2
		oauth2.RespondTo(resp).ErrInvalidClientBasicAuthorization(
			restSrv.serverCore.RealmName(), "")
		return
	}

	// To use this grant type, the client must be able to secure its credentials.
	if appIDNum := reqApp.RefKey.IDNum(); !appIDNum.IsService() && !appIDNum.IsUserAgentAuthorizationConfidential() {
		logReq(req.Request).
			Warn().Msgf("Client %v is not allowed to use grant type 'client_credentials'", reqApp.RefKey)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorUnauthorizedClient)
		return
	}

	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil && err != iam.ErrReqFieldAuthorizationTypeUnsupported {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Unable to read authorization")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsValid() {
		logCtx(reqCtx).
			Warn().Msg("Authorization context must not be valid")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	termDisplayName := ""

	regOutput := restSrv.serverCore.
		RegisterTerminal(iamserver.TerminalRegistrationInput{
			Context:        reqCtx,
			ApplicationRef: reqApp.RefKey,
			Data: iamserver.TerminalRegistrationInputData{
				UserRef:          ctxAuth.UserRef(),
				DisplayName:      termDisplayName,
				VerificationType: iam.TerminalVerificationResourceTypeOAuthClientCredentials,
				VerificationID:   0,
			}})
	if regOutput.Context.Err != nil {
		logCtx(reqCtx).
			Error().Err(regOutput.Context.Err).
			Msg("RegisterTerminal")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	termRef := regOutput.Data.TerminalRef
	termSecret := regOutput.Data.TerminalSecret
	issueTime := time.Now().UTC()

	accessToken, err := restSrv.serverCore.
		GenerateAccessTokenJWT(reqCtx, termRef, ctxAuth.UserRef(), issueTime)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).
			Msg("GenerateAccessTokenJWT")
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	//TODO: properly get the secret
	refreshToken, err := restSrv.serverCore.
		GenerateRefreshTokenJWT(reqCtx, termRef, termSecret, issueTime)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).
			Msg("GenerateRefreshTokenJWT")
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
			UserID:         ctxAuth.UserRef().AZIDText(),
			TerminalID:     termRef.AZIDText(),
			TerminalSecret: termSecret,
		})
}
