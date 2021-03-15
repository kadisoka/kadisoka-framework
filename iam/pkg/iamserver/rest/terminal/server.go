package terminal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"golang.org/x/text/language"

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
		Notes("The terminal registration endpoint is used to register "+
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
		Reads(iam.TerminalRegisterPostRequestJSONV1{}).
		Returns(http.StatusOK, "Terminal registered", iam.TerminalRegisterPostResponseJSONV1{}))

	restWS.Route(restWS.
		PUT("/fcm_registration_token").
		Metadata(restfulspec.KeyOpenAPITags, hidden).
		To(restSrv.putTerminalFCMRegistrationToken).
		Doc("Set terminal's FCM token").
		Notes("Associate the terminal with an FCM registration token. One token should "+
			"be associated to only one terminal.").
		Param(restWS.HeaderParameter("Authorization", sec.AuthorizationBearerAccessToken.String()).
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
			Warn().Err(err).Msg("Unable to read authorization")
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

	var terminalRegisterReq iam.TerminalRegisterPostRequestJSONV1
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

	if terminalRegisterReq.VerificationResourceType == "" {
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
		return
	}

	switch terminalRegisterReq.VerificationResourceType {
	case iam.TerminalVerificationResourceTypeEmailAddress:
		restSrv.handleTerminalRegisterByEmailAddress(
			resp, reqCtx, reqApp, terminalRegisterReq)
		return
	case iam.TerminalVerificationResourceTypePhoneNumber:
		restSrv.handleTerminalRegisterByPhoneNumber(
			resp, reqCtx, reqApp, terminalRegisterReq)
		return

	default:
		logCtx(reqCtx).
			Warn().Msgf("Unsupported verification resource type: %v", terminalRegisterReq.VerificationResourceType)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}
}

func (restSrv *Server) putTerminalFCMRegistrationToken(
	req *restful.Request, resp *restful.Response,
) {
	reqCtx, err := restSrv.RESTRequestContext(req.Request)
	if err != nil {
		logCtx(reqCtx).
			Err(err).Msg("Unable to read authorization")
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
	authClient *iam.Application,
	terminalRegisterReq iam.TerminalRegisterPostRequestJSONV1,
) {
	// Only for non-confidential user-agents
	if appRef := authClient.ID; !appRef.ID().IsUserAgentAuthorizationPublic() {
		logCtx(reqCtx).
			Warn().Msgf(
			"Client %v is not allowed to use this verification resource type", authClient.ID)
		rest.RespondTo(resp).EmptyError(
			http.StatusForbidden)
		return
	}

	termLangTags := restSrv.parseRequestAcceptLanguageTags(reqCtx, "")

	phoneNumber, err := iam.PhoneNumberFromString(terminalRegisterReq.VerificationResourceName)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Msgf(
			"Unable to parse verification resource name %s of type %s",
			terminalRegisterReq.VerificationResourceName,
			terminalRegisterReq.VerificationResourceType)
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

	termRef, _, codeExpiry, err := restSrv.serverCore.
		StartTerminalRegistrationByPhoneNumber(
			reqCtx, authClient.ID, phoneNumber,
			terminalRegisterReq.DisplayName, reqCtx.HTTPRequest().UserAgent(),
			termLangTags, verificationMethods)
	if err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).Msgf(
				"StartTerminalAuthorizationByPhoneNumber with %v failed", phoneNumber)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Err(err).Msgf(
			"StartTerminalAuthorizationByPhoneNumber with %v failed", phoneNumber)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	rest.RespondTo(resp).Success(
		&iam.TerminalRegisterPostResponseJSONV1{
			TerminalID: termRef.AZERText(),
			CodeExpiry: codeExpiry,
		})
	return
}

// terminal registration using email address
func (restSrv *Server) handleTerminalRegisterByEmailAddress(
	resp *restful.Response,
	reqCtx *iam.RESTRequestContext,
	authClient *iam.Application,
	terminalRegisterReq iam.TerminalRegisterPostRequestJSONV1,
) {
	if appRef := authClient.ID; !appRef.ID().IsUserAgentAuthorizationPublic() {
		logCtx(reqCtx).
			Warn().Msgf(
			"Client %v is not allowed to use this verification resource type", authClient.ID)
		rest.RespondTo(resp).EmptyError(
			http.StatusForbidden)
		return
	}

	termLangTags := restSrv.parseRequestAcceptLanguageTags(reqCtx, "")
	emailAddressStr := terminalRegisterReq.VerificationResourceName
	if !iam.IsValidEmailAddress(emailAddressStr) {
		logCtx(reqCtx).
			Warn().Msgf("Email address %v, is not valid", emailAddressStr)
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}
	emailAddress, err := iam.EmailAddressFromString(emailAddressStr)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Msgf(
			"Unable to parse verification resource name %s of type %s",
			terminalRegisterReq.VerificationResourceName,
			terminalRegisterReq.VerificationResourceType)
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

	termRef, _, codeExpiry, err := restSrv.serverCore.
		StartTerminalRegistrationByEmailAddress(
			reqCtx, authClient.ID, emailAddress,
			terminalRegisterReq.DisplayName, reqCtx.HTTPRequest().UserAgent(),
			termLangTags, verificationMethods)
	if err != nil {
		if errors.IsCallError(err) {
			logCtx(reqCtx).
				Warn().Err(err).Msgf("StartTerminalAuthorizationByEmailAddress with %v failed",
				emailAddress)
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		logCtx(reqCtx).
			Err(err).Msgf("StartTerminalAuthorizationByEmailAddress with %v failed",
			emailAddress)
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	rest.RespondTo(resp).Success(
		&iam.TerminalRegisterPostResponseJSONV1{
			TerminalID: termRef.AZERText(),
			CodeExpiry: codeExpiry,
		})
	return
}

func (restSrv *Server) handleTerminalRegisterByImplicit(
	resp *restful.Response,
	reqCtx *iam.RESTRequestContext,
	authClient *iam.Application,
	terminalRegisterReq iam.TerminalRegisterPostRequestJSONV1,
) {
	// Only if the client is able to secure its credentials.
	if !authClient.ID.ID().IsService() && !authClient.ID.ID().IsUserAgentAuthorizationConfidential() {
		logCtx(reqCtx).
			Warn().Msgf("Client %v is not allowed to use this verification resource type", authClient.ID)
		rest.RespondTo(resp).EmptyError(
			http.StatusForbidden)
		return
	}

	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsUserContext() {
		//TODO: determine if we should support user context
		logCtx(reqCtx).
			Warn().Msgf("Client %v is authenticating by implicit grant with valid user context", authClient.ID)
		rest.RespondTo(resp).EmptyError(
			http.StatusForbidden)
		return
	}

	//TODO: get from CallCtx
	termLangStrings := restSrv.parseRequestAcceptLanguage(reqCtx, "")
	termDisplayName := strings.TrimSpace(terminalRegisterReq.DisplayName)

	termRef, termSecret, err := restSrv.serverCore.
		RegisterTerminal(reqCtx, iamserver.TerminalRegistrationInput{
			ApplicationRef:   authClient.ID,
			UserRef:          iam.UserRefKeyZero(),
			DisplayName:      termDisplayName,
			AcceptLanguage:   strings.Join(termLangStrings, ","),
			VerificationType: terminalRegisterReq.VerificationResourceType,
			VerificationID:   0,
		})
	if err != nil {
		panic(err)
	}

	rest.RespondTo(resp).Success(
		&iam.TerminalRegisterPostResponseJSONV1{
			TerminalID:     termRef.AZERText(),
			TerminalSecret: termSecret,
		})
	return
}

// Parse accept languages from request
func (restSrv *Server) parseRequestAcceptLanguageTags(
	reqCtx *iam.RESTRequestContext,
	overrideAcceptLanguage string,
) (termLangTags []language.Tag) {
	termLangTags, _, err := language.ParseAcceptLanguage(overrideAcceptLanguage)
	if overrideAcceptLanguage != "" && err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Msgf("Unable to parse preferred languages from body %q", overrideAcceptLanguage)
	}
	if len(termLangTags) == 0 || err != nil {
		if httpReq := reqCtx.HTTPRequest(); httpReq != nil {
			var headerLangTags []language.Tag
			headerLangTags, _, err = language.
				ParseAcceptLanguage(httpReq.Header.Get("Accept-Language"))
			if err != nil {
				logCtx(reqCtx).
					Warn().Err(err).Msg("Unable to parse preferred languages from HTTP header")
			} else {
				if len(headerLangTags) > 0 {
					termLangTags = headerLangTags
				}
			}
		}
	}

	return termLangTags
}

// Parse accept languages from request
func (restSrv *Server) parseRequestAcceptLanguage(
	reqCtx *iam.RESTRequestContext,
	overrideAcceptLanguage string,
) (termLangStrings []string) {
	termLangTags := restSrv.parseRequestAcceptLanguageTags(reqCtx, overrideAcceptLanguage)
	for _, langTag := range termLangTags {
		termLangStrings = append(termLangStrings, langTag.String())
	}
	return termLangStrings
}

type terminalFCMRegistrationTokenPutRequest struct {
	Token string `json:"token"`
}
