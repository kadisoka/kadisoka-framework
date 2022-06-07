package user

import (
	"net/http"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/email"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
)

func (restSrv *Server) putUserEmailAddress(
	req *restful.Request, resp *restful.Response,
) {
	reqCtx, err := restSrv.RESTOpInputContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request context")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	var reqEntity UserEmailAddressPutRequestJSONV1
	err = req.ReadEntity(&reqEntity)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request entity")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	emailAddress := reqEntity.EmailAddress
	parsedEmailAddress, err := email.AddressFromString(emailAddress)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msgf("Email address %v, is not valid", emailAddress)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	var verificationMethods []eav10n.VerificationMethod
	for _, s := range reqEntity.VerificationMethods {
		m := eav10n.VerificationMethodFromString(s)
		if m.IsValid() {
			verificationMethods = append(verificationMethods, m)
		}
	}

	restSrv.handleSetEmailAddress(reqCtx, req, resp,
		parsedEmailAddress, verificationMethods)
}

func (restSrv *Server) handleSetEmailAddress(
	reqCtx *iam.RESTOpInputContext,
	req *restful.Request,
	resp *restful.Response,
	emailAddress email.Address,
	verificationMethods []eav10n.VerificationMethod,
) {
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsNotStaticallyValid() || !ctxAuth.IsUserSubject() {
		logCtx(reqCtx).
			Warn().Msgf("Unauthorized")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	if targetUserRefStr := req.PathParameter("user-id"); targetUserRefStr != "" && targetUserRefStr != "me" {
		targetUserRef, err := iam.UserRefKeyFromAZIDText(targetUserRefStr)
		if err != nil {
			logCtx(reqCtx).
				Warn().Err(err).
				Msg("Invalid user ID")
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		if !ctxAuth.IsUser(targetUserRef) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msg("Setting other user's email address is not allowed")
			rest.RespondTo(resp).EmptyError(
				http.StatusForbidden)
			return
		}
	}

	verificationID, codeExpiry, err := restSrv.serverCore.
		SetUserKeyEmailAddress(
			reqCtx, ctxAuth.UserRef(), emailAddress, verificationMethods)
	if err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("SetUserKeyEmailAddress %v", emailAddress)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("SetUserKeyEmailAddress %v", emailAddress)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	if verificationID == 0 {
		rest.RespondTo(resp).Success(nil)
		return
	}

	rest.RespondTo(resp).SuccessWithHTTPStatusCode(
		&UserEmailAddressPutResponse{
			VerificationID: verificationID,
			CodeExpiry:     *codeExpiry,
		},
		http.StatusAccepted)
}

//TODO(exa): should we allow confirming without the need to login
func (restSrv *Server) postUserEmailAddressVerificationConfirmation(
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

	var reqEntity UserEmailAddressVerificationConfirmationPostRequest
	err = req.ReadEntity(&reqEntity)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Request entity")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	updated, err := restSrv.serverCore.
		ConfirmUserEmailAddressVerification(
			reqCtx, reqEntity.VerificationID, reqEntity.Code)
	if err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("ConfirmUserEmailAddressVerification %v",
					reqEntity.VerificationID)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("ConfirmUserEmailAddressVerification %v",
				reqEntity.VerificationID)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	if !updated {
		rest.RespondTo(resp).EmptyError(http.StatusGone)
		return
	}

	rest.RespondTo(resp).Success(nil)
}

type UserEmailAddressPutRequestJSONV1 struct {
	EmailAddress        string   `json:"email_address"`
	VerificationMethods []string `json:"verification_methods"`
}

type UserEmailAddressPutResponse struct {
	VerificationID int64     `json:"verification_id"`
	CodeExpiry     time.Time `json:"code_expiry"`
}

type UserEmailAddressVerificationConfirmationPostRequest struct {
	VerificationID int64  `json:"verification_id"`
	Code           string `json:"code"`
}
