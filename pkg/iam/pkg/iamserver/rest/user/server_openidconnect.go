package user

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	oidc "github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/openid/connect"
	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/rest"
)

//TODO: the details would be depends on the type of the client:
// if it's internal, it could get all the details. Otherwise, it will
// be depended on the requested scope and user's privacy settings.
func (restSrv *Server) getUserOpenIDConnectUserInfo(
	req *restful.Request, resp *restful.Response,
) {
	reqCtx, err := restSrv.RESTCallInputContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request context")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsNotStaticallyValid() {
		logCtx(reqCtx).
			Warn().Msg("Unauthorized")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	requestedUserID := ctxAuth.UserID()

	userBaseProfile, err := restSrv.serverCore.
		GetUserBaseProfile(reqCtx, requestedUserID)
	if err != nil {
		logCtx(reqCtx).
			Err(err).Msg("User base profile fetch")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	userPhoneNumber, err := restSrv.serverCore.
		GetUserKeyPhoneNumber(reqCtx, requestedUserID)
	if err != nil {
		logCtx(reqCtx).
			Err(err).Msg("User phone number fetch")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	var phoneNumberStr string
	var phoneNumberVerified bool
	if userPhoneNumber != nil {
		phoneNumberStr = userPhoneNumber.String()
		phoneNumberVerified = true
	}

	//TODO(exa): should get display email address instead of primary
	// email address for this use case.
	userEmailAddress, err := restSrv.serverCore.
		GetUserKeyEmailAddress(reqCtx, requestedUserID)
	if err != nil {
		logCtx(reqCtx).
			Err(err).Msg("User email address fetch")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	var emailAddressStr string
	var emailAddressVerified bool
	if userEmailAddress != nil {
		emailAddressStr = userEmailAddress.RawInput()
		emailAddressVerified = true
	}

	userInfo := oidc.StandardClaims{
		Sub:                 requestedUserID.AZIDText(),
		Name:                userBaseProfile.DisplayName,
		Email:               emailAddressStr,
		EmailVerified:       emailAddressVerified,
		PhoneNumber:         phoneNumberStr,
		PhoneNumberVerified: phoneNumberVerified,
	}

	restSrv.eTagResponder.RespondGetJSON(req, resp, &userInfo)
}
