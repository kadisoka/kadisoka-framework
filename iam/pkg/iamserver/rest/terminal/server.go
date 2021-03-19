package terminal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/rest/logging"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/rest/sec"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
)

var (
	log    = logging.NewPkgLogger()
	logCtx = log.WithContext
	logReq = log.WithRequest
)

type ServerConfig struct {
	ServePath string
}

func NewServer(
	iamServerCore *iamserver.Core,
	config ServerConfig,
) *Server {
	return &Server{
		iamserver.RESTServiceServerWith(iamServerCore),
		config.ServePath,
	}
}

type Server struct {
	serverCore *iamserver.RESTServiceServerBase
	basePath   string
}

func (restSrv *Server) RESTRequestContext(req *http.Request) (*iam.RESTRequestContext, error) {
	return restSrv.serverCore.RESTRequestContext(req)
}

func (restSrv *Server) RestfulWebService() *restful.WebService {
	restWS := new(restful.WebService)
	restWS.
		Path(restSrv.basePath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"iam.v1.terminals"}
	hidden := append([]string{"hidden"}, tags...)

	restWS.Route(restWS.
		POST("/register").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		To(restSrv.postTerminalsRegister).
		Doc("Terminal registration endpoint").
		Notes(
			"The terminal registration endpoint is used to register "+
				"a terminal. This endpoint will send a verification code "+
				"through the configured external communication channel. "+
				"This code needs to be provided to the terminal secret "+
				"endpoint to obtain the secret of the terminal.\n\n"+
				"A **terminal** is a bound instance of client. It might or "+
				"might not be associated to a user.").
		Param(restWS.
			HeaderParameter(
				"Authorization", sec.AuthorizationBasicOAuth2ClientCredentials.String()).
			Required(true)).
		Reads(iam.TerminalRegistrationRequestJSONV1{}).
		Returns(http.StatusOK, "Terminal registered", iam.TerminalRegistrationResponseJSONV1{}))

	restWS.Route(restWS.
		DELETE("/{terminal-id}").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		To(restSrv.deleteTerminal).
		Doc("Terminal deletion (access-revocation) endpoint").
		Notes("This endpoint is used to revoke all access of a terminal.").
		Param(restWS.
			HeaderParameter(
				"Authorization", sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Param(restWS.
			PathParameter(
				"terminal-id", "The ID of the terminal to delete").
			Required(true)).
		Reads(iam.TerminalDeletionRequestJSONV1{}).
		Returns(http.StatusOK, "Terminal deleted", iam.TerminalDeletionResponseJSONV1{}))

	restWS.Route(restWS.
		PUT("/fcm_registration_token").
		Metadata(restfulspec.KeyOpenAPITags, hidden).
		To(restSrv.putTerminalFCMRegistrationToken).
		Doc("Set terminal's FCM token").
		Notes(
			"Associate the terminal with an FCM registration token. One "+
				"token should be associated to only one terminal.").
		Param(restWS.
			HeaderParameter(
				"Authorization", sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Reads(terminalFCMRegistrationTokenPutRequest{}).
		Returns(http.StatusNoContent, "Terminal's FCM token successfully set", nil))

	return restWS
}

func (restSrv *Server) postTerminalsRegister(
	req *restful.Request, resp *restful.Response,
) {
	reqApp, err := restSrv.serverCore.
		RequestApplication(req.Request)
	if err != nil {
		logReq(req.Request).
			Warn().Err(err).Msg("Client authentication")
		realmName := restSrv.serverCore.RealmName()
		if realmName == "" {
			realmName = "Restricted"
		}
		resp.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", realmName))
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	if reqApp == nil {
		logReq(req.Request).
			Warn().Msg("No authorized client")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil && err != iam.ErrReqFieldAuthorizationTypeUnsupported {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Unable to read authorization")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsValid() {
		logCtx(reqCtx).
			Warn().Msg("Authorization context must not be valid")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	var terminalRegisterReq iam.TerminalRegistrationRequestJSONV1
	err = req.ReadEntity(&terminalRegisterReq)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Msg("Unable to read entity from the request body")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	if terminalRegisterReq.VerificationResourceName == "" {
		logCtx(reqCtx).
			Warn().Msg("Resource name is missing")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	if iam.IsValidEmailAddress(terminalRegisterReq.VerificationResourceName) {
		restSrv.handleTerminalRegisterByEmailAddress(
			resp, reqCtx, reqApp, terminalRegisterReq)
		return
	}
	if _, err := iam.PhoneNumberFromString(terminalRegisterReq.VerificationResourceName); err == nil {
		restSrv.handleTerminalRegisterByPhoneNumber(
			resp, reqCtx, reqApp, terminalRegisterReq)
		return
	}

	logCtx(reqCtx).
		Warn().Msg("Resource verification type is missing")
	rest.RespondTo(resp).EmptyError(
		http.StatusBadRequest)
}

func (restSrv *Server) deleteTerminal(
	req *restful.Request, resp *restful.Response,
) {
	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Unable to read authorization")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	ctxAuth := reqCtx.Authorization()

	termIDStr := req.PathParameter("terminal-id")
	if termIDStr == "" {
		logCtx(reqCtx).
			Warn().Msg("Empty terminal ID")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	var termRef iam.TerminalRefKey
	if termIDStr == "self" {
		termRef = ctxAuth.TerminalRef()
	} else {
		termRef, err = iam.TerminalRefKeyFromAZERText(termIDStr)
		if err != nil {
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("Unabled to parse %q as a terminal ref-key", termIDStr)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
	}

	_, err = restSrv.serverCore.DeleteTerminal(reqCtx, termRef)
	if err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("DeleteTerminal with %v failed", termRef)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("DeleteTerminal with %v failed", termRef)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	rest.RespondTo(resp).
		Success(&iam.TerminalDeletionResponseJSONV1{})
}

func (restSrv *Server) putTerminalFCMRegistrationToken(
	req *restful.Request, resp *restful.Response,
) {
	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Unable to read authorization")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	if !reqCtx.IsUserContext() {
		logCtx(reqCtx).
			Warn().Msg("Unauthorized request")
		rest.RespondTo(resp).EmptyError(
			http.StatusUnauthorized)
		return
	}

	ctxAuth := reqCtx.Authorization()

	var fcmRegTokenReq terminalFCMRegistrationTokenPutRequest
	err = req.ReadEntity(&fcmRegTokenReq)
	if err != nil {
		panic(err)
	}

	if fcmRegTokenReq.Token == "" {
		logCtx(reqCtx).
			Warn().Msg("Empty token")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	err = restSrv.serverCore.
		SetTerminalFCMRegistrationToken(
			reqCtx, ctxAuth.TerminalRef(), ctxAuth.UserRef(),
			fcmRegTokenReq.Token)
	if err != nil {
		panic(err)
	}

	rest.RespondTo(resp).Success(nil)
}

// terminal register using phone number
func (restSrv *Server) handleTerminalRegisterByPhoneNumber(
	resp *restful.Response,
	reqCtx *iam.RESTRequestContext,
	authApp *iam.Application,
	terminalRegisterReq iam.TerminalRegistrationRequestJSONV1,
) {
	// Only for non-confidential user-agents
	if appRef := authApp.ID; !appRef.ID().IsUserAgentAuthorizationPublic() {
		logCtx(reqCtx).
			Warn().Msgf("Client %v is not allowed to use this verification resource type",
			authApp.ID)
		rest.RespondTo(resp).EmptyError(
			http.StatusForbidden)
		return
	}

	phoneNumber, err := iam.PhoneNumberFromString(terminalRegisterReq.VerificationResourceName)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Msgf("Unable to parse verification resource name %s as phone number",
			terminalRegisterReq.VerificationResourceName)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	var verificationMethods []pnv10n.VerificationMethod
	for _, s := range terminalRegisterReq.VerificationMethods {
		m := pnv10n.VerificationMethodFromString(s)
		if m != pnv10n.VerificationMethodUnspecified {
			verificationMethods = append(verificationMethods, m)
		}
	}

	authStartOutput := restSrv.serverCore.
		StartTerminalAuthorizationByPhoneNumber(
			iamserver.TerminalAuthorizationByPhoneNumberStartInput{
				Context:        reqCtx,
				ApplicationRef: authApp.ID,
				Data: iamserver.TerminalAuthorizationByPhoneNumberStartInputData{
					PhoneNumber:         phoneNumber,
					VerificationMethods: verificationMethods,
					TerminalAuthorizationStartInputBaseData: iamserver.TerminalAuthorizationStartInputBaseData{
						DisplayName: terminalRegisterReq.DisplayName,
					},
				},
			})
	if err = authStartOutput.Context.Err; err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("StartTerminalAuthorizationByPhoneNumber with %v failed", phoneNumber)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("StartTerminalAuthorizationByPhoneNumber with %v failed", phoneNumber)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	rest.RespondTo(resp).Success(
		&iam.TerminalRegistrationResponseJSONV1{
			TerminalID: authStartOutput.Data.TerminalRef.AZERText(),
			CodeExpiry: authStartOutput.Data.VerificationCodeExpiryTime,
		})
}

// terminal registration using email address
func (restSrv *Server) handleTerminalRegisterByEmailAddress(
	resp *restful.Response,
	reqCtx *iam.RESTRequestContext,
	authApp *iam.Application,
	terminalRegisterReq iam.TerminalRegistrationRequestJSONV1,
) {
	if appRef := authApp.ID; !appRef.ID().IsUserAgentAuthorizationPublic() {
		logCtx(reqCtx).
			Warn().Msgf("Client %v is not allowed to use this verification resource type",
			authApp.ID)
		rest.RespondTo(resp).EmptyError(
			http.StatusForbidden)
		return
	}

	emailAddressStr := terminalRegisterReq.VerificationResourceName
	if !iam.IsValidEmailAddress(emailAddressStr) {
		logCtx(reqCtx).
			Warn().Msgf("Provided email address %v is not valid", emailAddressStr)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}
	emailAddress, err := iam.EmailAddressFromString(emailAddressStr)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msgf("Unable to parse %s as email address",
				terminalRegisterReq.VerificationResourceName)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	var verificationMethods []eav10n.VerificationMethod
	for _, s := range terminalRegisterReq.VerificationMethods {
		m := eav10n.VerificationMethodFromString(s)
		if m != eav10n.VerificationMethodUnspecified {
			verificationMethods = append(verificationMethods, m)
		}
	}

	authStartOutput := restSrv.serverCore.
		StartTerminalAuthorizationByEmailAddress(
			iamserver.TerminalAuthorizationByEmailAddressStartInput{
				Context:        reqCtx,
				ApplicationRef: authApp.ID,
				Data: iamserver.TerminalAuthorizationByEmailAddressStartInputData{
					EmailAddress:        emailAddress,
					VerificationMethods: verificationMethods,
					TerminalAuthorizationStartInputBaseData: iamserver.TerminalAuthorizationStartInputBaseData{
						DisplayName: terminalRegisterReq.DisplayName,
					},
				},
			})
	if err = authStartOutput.Context.Err; err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("StartTerminalAuthorizationByEmailAddress with %v failed", emailAddress)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("StartTerminalAuthorizationByEmailAddress with %v failed", emailAddress)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	rest.RespondTo(resp).Success(
		&iam.TerminalRegistrationResponseJSONV1{
			TerminalID: authStartOutput.Data.TerminalRef.AZERText(),
			CodeExpiry: authStartOutput.Data.VerificationCodeExpiryTime,
		})
}

func (restSrv *Server) handleTerminalRegisterByImplicit(
	resp *restful.Response,
	reqCtx *iam.RESTRequestContext,
	authApp *iam.Application,
	terminalRegisterReq iam.TerminalRegistrationRequestJSONV1,
) {
	// Only if the client is able to secure its credentials.
	if !authApp.ID.ID().IsService() && !authApp.ID.ID().IsUserAgentAuthorizationConfidential() {
		logCtx(reqCtx).
			Warn().Msgf("Client %v is not allowed to use this verification resource type", authApp.ID)
		rest.RespondTo(resp).EmptyError(
			http.StatusForbidden)
		return
	}

	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsUserContext() {
		//TODO: determine if we should support user context
		logCtx(reqCtx).
			Warn().Msgf("Client %v is authenticating by implicit grant with valid user context", authApp.ID)
		rest.RespondTo(resp).EmptyError(
			http.StatusForbidden)
		return
	}

	termDisplayName := strings.TrimSpace(terminalRegisterReq.DisplayName)

	regOutput := restSrv.serverCore.
		RegisterTerminal(iamserver.TerminalRegistrationInput{
			Context:        reqCtx,
			ApplicationRef: authApp.ID,
			Data: iamserver.TerminalRegistrationInputData{
				UserRef:          iam.UserRefKeyZero(),
				DisplayName:      termDisplayName,
				VerificationType: iam.TerminalVerificationResourceTypeOAuthImplicit,
				VerificationID:   0,
			}})
	if regOutput.Context.Err != nil {
		panic(regOutput.Context.Err)
	}

	rest.RespondTo(resp).Success(
		&iam.TerminalRegistrationResponseJSONV1{
			TerminalID:     regOutput.Data.TerminalRef.AZERText(),
			TerminalSecret: regOutput.Data.TerminalSecret,
		})
}

type terminalFCMRegistrationTokenPutRequest struct {
	Token string `json:"token"`
}
