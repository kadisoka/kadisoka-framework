package oauth2

import (
	"net/http"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"golang.org/x/text/language"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/oauth2"
	apperrs "github.com/kadisoka/kadisoka-framework/foundation/pkg/app/errors"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/rest/logging"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam/rest/sec"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
)

var (
	log    = logging.NewPkgLogger()
	logCtx = log.WithContext
	logReq = log.WithRequest
)

type ServerConfig struct {
	ServePath string
	SignInURL string
}

// New instantiates an Server.
func NewServer(
	iamServerCore *iamserver.Core,
	config ServerConfig,
) (*Server, error) {
	if !iamServerCore.JWTKeyChain().CanSign() {
		return nil, apperrs.NewConfigurationMsg("JWT key chain is required")
	}
	return &Server{
		iamserver.RESTServiceServerWith(iamServerCore),
		config.ServePath,
		config.SignInURL,
	}, nil
}

// Server is a limited implementation of OAuth 2.0 Authorization Framework standard (RFC 6749)
type Server struct {
	serverCore *iamserver.RESTServiceServerBase
	basePath   string
	signInURL  string
}

func (restSrv *Server) jwtKeyChain() *iam.JWTKeyChain {
	return restSrv.serverCore.JWTKeyChain()
}

func (restSrv *Server) RESTRequestContext(req *http.Request) (*iam.RESTRequestContext, error) {
	return restSrv.serverCore.RESTRequestContext(req)
}

// RestfulWebService is used to obtain restful WebService with all endpoints set up.
func (restSrv *Server) RestfulWebService() *restful.WebService {
	restWS := new(restful.WebService)
	restWS.
		Path(restSrv.basePath).
		Consumes("application/x-www-form-urlencoded").
		Produces(restful.MIME_JSON)

	tags := []string{"iam.v1.oauth2"}

	restWS.Route(restWS.
		GET("/authorize").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		To(restSrv.getAuthorize).
		Doc("OAuth 2.0 authorization endpoint").
		Notes(
			"This endpoint is the standard-conforming endpoint.\n\nThis "+
				"endpoint is used by client/consumer applications to request "+
				"authorization for any of the users.").
		Param(restWS.
			QueryParameter(
				"client_id", "The ID of the client which makes the request").
			Required(true)).
		Param(restWS.
			QueryParameter(
				"response_type", "Value MUST be set to `code`").
			Required(true)).
		Param(restWS.
			QueryParameter(
				"redirect_uri", "Any of client's registered redirection URIs")).
		Param(restWS.
			QueryParameter(
				"state", "An opaque value used by the client to "+
					"maintain state between the request and callback.")).
		Returns(http.StatusFound, "Success", nil))

	restWS.Route(restWS.
		POST("/authorize").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		To(restSrv.postAuthorize).
		Doc("Authorization endpoint").
		Notes(
			"This endpoint is not defined in the standard.\n\nThis endpoint "+
				"is used by the web front-end when a resource owner granted "+
				"the authorization. All the parameters are mirroring the "+
				"standard endpoint except that this endpoint requires "+
				"bearer access token as the value of Authorization header.").
		Param(restWS.
			HeaderParameter(
				"Authorization", sec.AuthorizationBearerAccessToken.String()).
			Required(true)).
		Param(restWS.
			FormParameter(
				"client_id", "The ID of the client which makes the request").
			Required(true)).
		Param(restWS.
			FormParameter(
				"response_type", "Value MUST be set to `code`").
			Required(true)).
		Param(restWS.
			FormParameter(
				"redirect_uri", "Any of client's registered redirection URIs")).
		Param(restWS.
			FormParameter(
				"state", "An opaque value used by the client to "+
					"maintain state between the request and callback.")).
		Returns(http.StatusOK, "Success", iam.OAuth2AuthorizePostResponse{}))

	restWS.Route(restWS.
		GET("/jwks").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		To(restSrv.getJWKS).
		Doc("JSON Web Key Set endpoint").
		Notes("The JSON Web Key Set endpoint provides public keys needed "+
			"to verify JWT (JSON Web Token) tokens issued by this service.").
		Returns(http.StatusOK, "OK. See https://tools.ietf.org/html/rfc7517 for the data structure", nil))

	restWS.Route(restWS.
		POST("/token").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		To(restSrv.postToken).
		Doc("OAuth token endpoint").
		Notes(
			"The token endpoint is used by the client to obtain an "+
				"access token by presenting its authorization grant or "+
				"refresh token. The token endpoint is used with every "+
				"authorization grant except for the implicit grant type "+
				"(since an access token is issued directly). RFC 6749 § 3.2.").
		Param(restWS.
			HeaderParameter(
				"Authorization", sec.AuthorizationBasicOAuth2ClientCredentials.String()).
			Required(true)).
		Param(restWS.
			FormParameter(
				"grant_type", "Supported grant types: `password`, "+
					"`authorization_code`, `client_credentials`").
			Required(true)).
		Param(restWS.
			FormParameter(
				"username", "Required for `password` grant type")).
		Param(restWS.
			FormParameter(
				"password", "For use with `password` grant type")).
		Param(restWS.
			FormParameter(
				"code", "Required for `authorization_code` grant type")).
		Returns(http.StatusOK, "Authorization successful", iam.OAuth2TokenResponse{}).
		Returns(http.StatusBadRequest, "Request has missing data or contains invalid data", oauth2.ErrorResponse{}).
		Returns(http.StatusUnauthorized, "Client authorization check failure", oauth2.ErrorResponse{}))

	return restWS
}

// Parse preferred languages from request
func (restSrv *Server) parseRequestAcceptLanguage(
	req *restful.Request,
	reqCtx *iam.RESTRequestContext,
) (langTagSet []language.Tag) {
	langTagSet, _, err := language.ParseAcceptLanguage(
		req.Request.Header.Get("Accept-Language"))
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msg("Unable to parse preferred languages from HTTP header")
	}

	return langTagSet
}
