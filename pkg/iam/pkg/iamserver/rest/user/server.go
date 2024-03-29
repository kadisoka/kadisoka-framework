package user

import (
	"net/http"

	restfulopenapi "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	oidc "github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/openid/connect"
	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam/rest/logging"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam/rest/sec"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver"
)

var (
	log    = logging.NewPkgLogger()
	logCtx = log.WithContext
)

type ServerConfig struct {
	ServePath string
}

func NewServer(
	iamServerCore *iamserver.Core,
	config ServerConfig,
) *Server {
	return &Server{
		serverCore:    iamserver.RESTServiceServerWith(iamServerCore),
		basePath:      config.ServePath,
		eTagResponder: rest.NewETagResponder(512),
	}
}

type Server struct {
	serverCore    *iamserver.RESTServiceServerBase
	basePath      string
	eTagResponder *rest.ETagResponder
}

func (restSrv *Server) RESTCallInputContext(req *http.Request) (*iam.RESTCallInputContext, error) {
	return restSrv.serverCore.RESTCallInputContext(req)
}

func (restSrv *Server) RestfulWebService() *restful.WebService {
	restWS := new(restful.WebService)
	restWS.
		Path(restSrv.basePath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"iam.v1.users"}
	hidden := append([]string{"hidden"}, tags...)

	restWS.Route(restWS.
		GET("/{user-id}").
		To(restSrv.getUser).
		Metadata(restfulopenapi.KeyOpenAPITags, hidden).
		Doc("Retrieve basic profile of current user").
		Produces(restful.MIME_JSON, "application/x-protobuf").
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Param(restWS.PathParameter("user-id",
			"Set to a valid user ID or 'me'.").
			Required(true)).
		Returns(http.StatusOK, "OK", iam.UserJSONV1{}))

	restWS.Route(restWS.
		PUT("/me/password").
		To(restSrv.putUserPassword).
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		Doc("Set password for registered users").
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Reads(userPasswordPutRequest{}).
		Returns(http.StatusBadRequest, "Request has missing data or contains invalid data", rest.ErrorResponse{}).
		Returns(http.StatusUnauthorized, "Client authorization check failure", rest.ErrorResponse{}).
		Returns(http.StatusConflict, "Request has duplicate value or contains invalid data", rest.ErrorResponse{}).
		Returns(http.StatusNoContent, "Password set", nil))

	restWS.Route(restWS.
		PUT("/{user-id}/email_address").
		To(restSrv.putUserEmailAddress).
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		Doc("Set a new login email address for the current user").
		Notes("The email address needs to be verified before it's set as user's login "+
			"email address.").
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Param(restWS.PathParameter("user-id", "The ID of the user or `me`")).
		Reads(UserEmailAddressPutRequestJSONV1{}).
		Returns(http.StatusAccepted,
			"Email address is accepted by the server and waiting for verification confirmation",
			&UserEmailAddressPutResponse{}).
		Returns(http.StatusNoContent,
			"Provided email address is same as current one.",
			nil))

	restWS.Route(restWS.
		POST("/me/email_address/verification_confirmation").
		To(restSrv.postUserEmailAddressVerificationConfirmation).
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		Doc("Confirm email address verification").
		Reads(UserEmailAddressVerificationConfirmationPostRequest{}).
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				sec.AuthorizationBearerAccessToken.String()).
			Required(false)).
		Returns(http.StatusNoContent,
			"User login email address successfully set", nil))

	restWS.Route(restWS.
		PUT("/me/phone_number").
		To(restSrv.putUserPhoneNumber).
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		Doc("Set a new login phone number for the current user").
		Notes("The phone number needs to be verified before it's set as user's login "+
			"phone number.").
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Reads(UserPhoneNumberPutRequest{}).
		Returns(
			http.StatusAccepted,
			"Phone number is accepted by the server and waiting for verification confirmation",
			&UserPhoneNumberPutResponse{}).
		Returns(
			http.StatusNoContent,
			"Provided phone number is same as current one.",
			nil))

	restWS.Route(restWS.
		POST("/me/phone_number/verification_confirmation").
		To(restSrv.postUserPhoneNumberVerificationConfirmation).
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		Doc("Confirm phone number verification").
		Reads(UserPhoneNumberVerificationConfirmationPostRequest{}).
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				sec.AuthorizationBearerAccessToken.String()).
			Required(false)).
		Returns(
			http.StatusNoContent,
			"User login phone number successfully set", nil))

	restWS.Route(restWS.
		PUT("/me/profile_image").
		Consumes("multipart/form-data").
		To(restSrv.putUserProfileImage).
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		Doc("Update user profile image").
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Param(restWS.
			FormParameter(
				"body", "File to upload").
			DataType("file").
			Required(true)).
		Returns(http.StatusInternalServerError, "An unexpected condition was encountered in processing the request", nil).
		Returns(http.StatusBadRequest, "The server cannot or will not process the request due to an apparent client error", nil).
		Returns(http.StatusUnauthorized, "Authentication is required and has failed or has not yet been provided", nil).
		Returns(http.StatusNotAcceptable, "The target resource does not have a current representation that would be acceptable.", nil).
		Returns(http.StatusOK, "Profile image updated", userProfileImagePutResponse{}))

	restWS.Route(restWS.
		GET("/me/openidconnect-userinfo").
		To(restSrv.getUserOpenIDConnectUserInfo).
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		Doc("Retrieve Claims about the authenticated End-User").
		Notes("The UserInfo Endpoint is an OAuth 2.0 Protected "+
			"Resource that returns Claims about the authenticated "+
			"End-User. To obtain the requested Claims about the End-User, "+
			"the Client makes a request to the UserInfo Endpoint using an "+
			"Access Token obtained through OpenID Connect Authentication. "+
			"These Claims are represented by a JSON object that contains a "+
			"collection of name and value pairs for the Claims.").
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Returns(http.StatusOK, "OK", oidc.StandardClaims{}))

	return restWS
}

func (restSrv *Server) getUser(req *restful.Request, resp *restful.Response) {
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

	requestedUserIDStr := req.PathParameter("user-id")
	if requestedUserIDStr == "" {
		logCtx(reqCtx).
			Warn().Msg("Invalid parameter value path.user-id: empty")
		rest.RespondTo(resp).EmptyError(
			http.StatusBadRequest)
		return
	}

	var requestedUserID iam.UserID
	if requestedUserIDStr == "me" {
		if !reqCtx.IsUserContext() {
			logCtx(reqCtx).
				Warn().Msg("Invalid request: 'me' can only be used with user access token")
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
		requestedUserID = ctxAuth.UserID()
	} else {
		requestedUserID, err = iam.UserIDFromAZIDText(requestedUserIDStr)
		if err != nil {
			logCtx(reqCtx).
				Warn().Err(err).Msg("Invalid parameter value path.user-id")
			rest.RespondTo(resp).EmptyError(
				http.StatusBadRequest)
			return
		}
	}

	if acceptType := req.Request.Header.Get("Accept"); acceptType == "application/x-protobuf" {
		userInfo, err := restSrv.serverCore.
			GetUserInfoV1(reqCtx, requestedUserID)
		if err != nil {
			panic(err)
		}
		restSrv.eTagResponder.RespondGetProtoMessage(req, resp, userInfo)
		return
	}

	userBaseProfile, err := restSrv.serverCore.
		GetUserBaseProfile(reqCtx, requestedUserID)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).
			Msg("User base profile fetch")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}

	restUserProfile := iam.UserJSONV1FromBaseProfile(userBaseProfile, requestedUserID)

	userPhoneNumber, err := restSrv.serverCore.
		GetUserKeyPhoneNumber(reqCtx, requestedUserID)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).
			Msg("User phone number fetch")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	if userPhoneNumber != nil {
		restUserProfile.Data.PhoneNumber = userPhoneNumber.String()
	}

	//TODO(exa): should get display email address instead of primary
	// email address for this use case.
	userEmailAddress, err := restSrv.serverCore.
		GetUserKeyEmailAddress(reqCtx, requestedUserID)
	if err != nil {
		logCtx(reqCtx).
			Error().Err(err).
			Msg("User email address fetch")
		rest.RespondTo(resp).EmptyError(
			http.StatusInternalServerError)
		return
	}
	if userEmailAddress != nil {
		restUserProfile.Data.EmailAddress = userEmailAddress.RawInput()
	}

	restSrv.eTagResponder.RespondGetJSON(req, resp, restUserProfile)
}
