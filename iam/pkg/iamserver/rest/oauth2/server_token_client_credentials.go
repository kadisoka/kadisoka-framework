//

package oauth2

import (
	"github.com/emicklei/go-restful/v3"
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
	if appIDNum := reqApp.ID.IDNum(); !appIDNum.IsService() && !appIDNum.IsUserAgentAuthorizationConfidential() {
		logReq(req.Request).
			Warn().Msgf("Client %v is not allowed to use grant type 'client_credentials'", reqApp.ID)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorUnauthorizedClient)
		return
	}

	reqCtx, err := restSrv.RESTCallInputContext(req.Request)
	if err != nil && err != iam.ErrReqFieldAuthorizationTypeUnsupported {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Unable to read authorization")
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

	termDisplayName := ""

	regOutput := restSrv.serverCore.
		RegisterTerminal(iamserver.TerminalRegistrationInput{
			Context:       reqCtx,
			ApplicationID: reqApp.ID,
			Data: iamserver.TerminalRegistrationInputData{
				UserID:           ctxAuth.UserID(),
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

	termID := regOutput.Data.TerminalID
	termSecret := regOutput.Data.TerminalSecret

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
			UserID:         ctxAuth.UserID().AZIDText(),
			TerminalID:     termID.AZIDText(),
			TerminalSecret: termSecret,
		})
}
