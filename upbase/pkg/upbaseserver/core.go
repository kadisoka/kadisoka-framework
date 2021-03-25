package upbaseserver

import (
	"net/http"

	"github.com/alloyzeus/go-azfl/azfl/errors"

	oidc "github.com/kadisoka/kadisoka-framework/foundation/pkg/api/openid/connect"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/upbase/pkg/upbase"
)

type Core struct {
	iamSvc iam.ConsumerServer
}

func (srvCore *Core) RESTRequestContext(
	req *http.Request,
) (*upbase.RESTRequestContext, error) {
	iamReqCtx, err := srvCore.iamSvc.RESTRequestContext(req)
	return &upbase.RESTRequestContext{RESTRequestContext: *iamReqCtx}, err
}

func (srvCore *Core) GetUserOpenIDConnectStandardClaims(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*oidc.StandardClaims, error) {
	if callCtx == nil {
		return nil, iam.ErrOperationContextMissing
	}

	ctxAuth := callCtx.Authorization()
	if !ctxAuth.IsUser(userRef) {
		return nil, iam.ErrOperationNotAllowed
	}
	//TODO(exa): ensure that the context user has the privilege

	userInstInfo, err := srvCore.iamSvc.GetUserInstanceInfo(callCtx, userRef)
	if err != nil {
		//TODO: translate error
		return nil, errors.Wrap("GetUserInstanceInfo", err)
	}

	if userInstInfo == nil {
		return nil, nil
	}

	return srvCore.getUserOpenIDConnectStandardClaimsNoAC(callCtx, userRef)
}

func (srvCore *Core) getUserOpenIDConnectStandardClaimsNoAC(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) (*oidc.StandardClaims, error) {
	return nil, errors.ErrUnimplemented

	// userBaseProfile, err := restSrv.serverCore.
	// 	GetUserBaseProfile(reqCtx, requestedUserRef)
	// if err != nil {
	// 	logCtx(reqCtx).
	// 		Err(err).Msg("User base profile fetch")
	// 	rest.RespondTo(resp).EmptyError(
	// 		http.StatusInternalServerError)
	// 	return
	// }

	// userPhoneNumber, err := restSrv.serverCore.
	// 	GetUserKeyPhoneNumber(reqCtx, requestedUserRef)
	// if err != nil {
	// 	logCtx(reqCtx).
	// 		Err(err).Msg("User phone number fetch")
	// 	rest.RespondTo(resp).EmptyError(
	// 		http.StatusInternalServerError)
	// 	return
	// }
	// var phoneNumberStr string
	// var phoneNumberVerified bool
	// if userPhoneNumber != nil {
	// 	phoneNumberStr = userPhoneNumber.String()
	// 	phoneNumberVerified = true
	// }

	// //TODO(exa): should get display email address instead of primary
	// // email address for this use case.
	// userEmailAddress, err := restSrv.serverCore.
	// 	GetUserKeyEmailAddress(reqCtx, requestedUserRef)
	// if err != nil {
	// 	logCtx(reqCtx).
	// 		Err(err).Msg("User email address fetch")
	// 	rest.RespondTo(resp).EmptyError(
	// 		http.StatusInternalServerError)
	// 	return
	// }
	// var emailAddressStr string
	// var emailAddressVerified bool
	// if userEmailAddress != nil {
	// 	emailAddressStr = userEmailAddress.RawInput()
	// 	emailAddressVerified = true
	// }

	// userInfo := oidc.StandardClaims{
	// 	Sub:                 userRef.AZIDText(),
	// 	Name:                userBaseProfile.DisplayName,
	// 	Email:               emailAddressStr,
	// 	EmailVerified:       emailAddressVerified,
	// 	PhoneNumber:         phoneNumberStr,
	// 	PhoneNumberVerified: phoneNumberVerified,
	// }

	// return &userInfo, nil
}
