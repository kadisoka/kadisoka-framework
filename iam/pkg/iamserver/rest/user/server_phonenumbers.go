package user

import (
	"net/http"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/emicklei/go-restful"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

func (restSrv *Server) putUserPhoneNumber(
	req *restful.Request, resp *restful.Response,
) {
	reqCtx, err := restSrv.RESTOpInputContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request context")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsNotStaticallyValid() && !ctxAuth.IsUserSubject() {
		logCtx(reqCtx).
			Warn().Msg("Unauthorized")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	var reqEntity UserPhoneNumberPutRequest
	err = req.ReadEntity(&reqEntity)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Body read")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	phoneNumber, err := telephony.PhoneNumberFromString(reqEntity.PhoneNumber)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msgf("Unable to parse %q as phone number",
				reqEntity.PhoneNumber)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}
	if !phoneNumber.IsSound() {
		logCtx(reqCtx).
			Warn().Msgf("Provided phone number %q is not valid", reqEntity.PhoneNumber)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	var verificationMethods []pnv10n.VerificationMethod
	for _, s := range reqEntity.VerificationMethods {
		m := pnv10n.VerificationMethodFromString(s)
		if m.IsValid() {
			verificationMethods = append(verificationMethods, m)
		}
	}

	verificationID, codeExpiry, err := restSrv.serverCore.
		SetUserKeyPhoneNumber(
			reqCtx, ctxAuth.UserRef(), phoneNumber, verificationMethods)
	if err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("SetUserKeyPhoneNumber to %v",
					phoneNumber)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("SetUserKeyPhoneNumber to %v",
				phoneNumber)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	if verificationID == 0 {
		rest.RespondTo(resp).Success(nil)
		return
	}

	rest.RespondTo(resp).SuccessWithHTTPStatusCode(
		&UserPhoneNumberPutResponse{
			VerificationID: verificationID,
			CodeExpiry:     *codeExpiry,
		},
		http.StatusAccepted)
	return
}

//TODO(exa): should we allow confirming without the need to login
func (restSrv *Server) postUserPhoneNumberVerificationConfirmation(
	req *restful.Request, resp *restful.Response,
) {
	reqCtx, err := restSrv.RESTOpInputContext(req.Request)
	if !reqCtx.IsUserContext() {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request context")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	var reqEntity UserPhoneNumberVerificationConfirmationPostRequest
	err = req.ReadEntity(&reqEntity)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Unable to load request content")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	updated, err := restSrv.serverCore.
		ConfirmUserPhoneNumberVerification(
			reqCtx, reqEntity.VerificationID, reqEntity.Code)
	if err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("ConfirmUserPhoneNumberVerification %v failed",
					reqEntity.VerificationID)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("ConfirmUserPhoneNumberVerification %v failed",
				reqEntity.VerificationID)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	if !updated {
		rest.RespondTo(resp).EmptyError(
			http.StatusGone)
		return
	}

	rest.RespondTo(resp).Success(nil)
}

type UserPhoneNumberPutRequest struct {
	PhoneNumber         string   `json:"phone_number"`
	VerificationMethods []string `json:"verification_methods"`
}

type UserPhoneNumberPutResponse struct {
	VerificationID int64     `json:"verification_id"`
	CodeExpiry     time.Time `json:"code_expiry"`
}

type UserPhoneNumberVerificationConfirmationPostRequest struct {
	VerificationID int64  `json:"verification_id"`
	Code           string `json:"code"`
}
