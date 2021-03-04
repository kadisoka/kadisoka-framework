//

package oauth2

import (
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/oauth2"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
)

func (restSrv *Server) handleTokenRequestByClientCredentials(
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

	// To use this grant type, the client must be able to secure its credentials.
	if !reqClient.ID.IsConfidential() {
		logReq(req.Request).
			Warn().Msgf("Client %v is not allowed to use grant type 'client_credentials'", reqClient.ID)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorUnauthorizedClient)
		return
	}

	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil && err != iam.ErrReqFieldAuthorizationTypeUnsupported {
		logCtx(reqCtx).
			Warn().Err(err).Msg("Unable to read authorization")
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

	preferredLanguages := restSrv.parseRequestAcceptLanguage(req, reqCtx, "")
	termDisplayName := ""

	termID, termSecret, err := restSrv.serverCore.
		RegisterTerminal(reqCtx, iamserver.TerminalRegistrationInput{
			ClientID:         reqClient.ID,
			UserRef:          authCtx.UserRef(),
			DisplayName:      termDisplayName,
			AcceptLanguage:   strings.Join(preferredLanguages, ","),
			VerificationType: iam.TerminalVerificationResourceTypeOAuthClientCredentials,
			VerificationID:   0,
		})
	if err != nil {
		logCtx(reqCtx).
			Error().Msgf("RegisterTerminal: %v", err)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	issueTime := time.Now().UTC()

	accessToken, err := restSrv.serverCore.
		GenerateAccessTokenJWT(reqCtx, termID, authCtx.UserRef(), issueTime)
	if err != nil {
		logCtx(reqCtx).
			Error().Msgf("GenerateAccessTokenJWT: %v", err)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}

	//TODO: properly get the secret
	refreshToken, err := restSrv.serverCore.
		GenerateRefreshTokenJWT(reqCtx, termID, termSecret, issueTime)
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
			UserID:         authCtx.UserRef().AZERText(),
			TerminalID:     termID.String(),
			TerminalSecret: termSecret,
		})
}
