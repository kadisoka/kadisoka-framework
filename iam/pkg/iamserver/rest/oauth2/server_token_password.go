//

package oauth2

import (
	"strings"
	"time"

	"github.com/emicklei/go-restful"
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
			Warn().Err(err).Msg("Client authentication")
		// RFC 6749 § 5.2
		oauth2.RespondTo(resp).ErrInvalidClientBasicAuthorization(
			restSrv.serverCore.RealmName(), "")
		return
	}

	if reqApp != nil && !reqApp.ID.ID().IsUserAgentAuthorizationConfidential() {
		logReq(req.Request).
			Warn().Msgf("Client %v is not authorized to use grant type password", reqApp.ID)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorUnauthorizedClient)
		return
	}

	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil && err != iam.ErrReqFieldAuthorizationTypeUnsupported {
		logCtx(reqCtx).
			Warn().Msgf("Unable to read authorization: %v", err)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}
	authCtx := reqCtx.Authorization()
	if authCtx.IsValid() {
		logCtx(reqCtx).
			Warn().Msgf("Authorization context must not be valid")
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
				Warn().Msgf("Unrecognized identifier scheme: %v", names[0])
		}
	}

	logReq(req.Request).
		Warn().Msgf("Password grant with no scheme.")
	oauth2.RespondTo(resp).ErrorCode(
		oauth2.ErrorInvalidGrant)
}

func (restSrv *Server) handleTokenRequestByPasswordGrantWithTerminalCreds(
	reqCtx *iam.RESTRequestContext,
	resp *restful.Response,
	reqApp *iam.Application,
	terminalIDStr string,
	terminalSecret string,
) {
	termRef, err := iam.TerminalRefKeyFromAZERText(terminalIDStr)
	if err != nil {
		logCtx(reqCtx).
			Warn().Msgf("Unable to parse username %q as TerminalID: %v", terminalIDStr, err)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	authOK, userRef, err := restSrv.serverCore.
		AuthenticateTerminal(termRef, terminalSecret)
	if err != nil {
		logCtx(reqCtx).
			Error().Msgf("Terminal %v authentication failed: %v", termRef, err)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorServerError)
		return
	}
	if !authOK {
		logCtx(reqCtx).
			Warn().Msgf("Terminal %v authentication failed", termRef)
		oauth2.RespondTo(resp).ErrorCode(
			oauth2.ErrorInvalidGrant)
		return
	}

	if userRef.IsValid() {
		userInstInfo, err := restSrv.serverCore.
			GetUserInstanceInfo(userRef)
		if err != nil {
			logCtx(reqCtx).
				Warn().Msgf("Terminal %v user account state: %v", termRef, err)
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorServerError)
			return
		}
		if userInstInfo == nil || !userInstInfo.IsActive() {
			var status string
			if userInstInfo == nil {
				status = "not exist"
			} else {
				status = "deleted"
			}
			logCtx(reqCtx).
				Warn().Msgf("User %v %s", userRef, status)
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorInvalidGrant)
			return
		}
	}

	if reqApp != nil {
		if !reqApp.ID.EqualsApplicationRefKey(termRef.Application()) {
			logCtx(reqCtx).
				Error().Msgf("Terminal %v is not associated to client %v", termRef, reqApp.ID)
			oauth2.RespondTo(resp).ErrorCode(
				oauth2.ErrorServerError)
			return
		}
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
			UserID: userRef.AZERText(),
		})
}
