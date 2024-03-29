package main

import (
	"net/http"

	restfulopenapi "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

func NewRESTService(
	iamConsumerServer iam.ConsumerServer,
	basePath string,
) *RESTService {
	return &RESTService{
		iamCS:    iamConsumerServer,
		basePath: basePath,
	}
}

type RESTService struct {
	iamCS    iam.ConsumerServer
	basePath string
}

func (restSvc *RESTService) RestfulWebService() *restful.WebService {
	restWS := new(restful.WebService)
	restWS.Path(restSvc.basePath).
		Produces(restful.MIME_JSON)

	tags := []string{"Microservice"}

	restWS.Route(restWS.
		GET("/auth").
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		To(restSvc.getAuth).
		Doc("Obtain access token using parameters obtained from a OAuth 2.0 authorization code flow").
		Param(restWS.
			QueryParameter(
				"authorization_code",
				"The authorization code received from the authorization server.").
			Required(true)).
		Param(restWS.
			QueryParameter(
				"state",
				"Will be provided by authorization server if the `state` "+
					"parameter was present in the client authorization request.")).
		Returns(http.StatusOK, "OK", &authGetResponse{}))

	restWS.Route(restWS.
		GET("/hello").
		Metadata(restfulopenapi.KeyOpenAPITags, tags).
		To(restSvc.getHello).
		Doc("Hello").
		Param(restWS.
			HeaderParameter(
				iam.AuthorizationMetadataKey,
				"Bearer access_token").
			Required(true)).
		Returns(http.StatusOK, "OK", &helloGetResponse{}))

	return restWS
}

type authGetResponse struct {
	AccessToken string `json:"access_token"`
}

func (restSvc *RESTService) getAuth(req *restful.Request, resp *restful.Response) {
	reqCtx, err := restSvc.iamCS.RESTCallInputContext(req.Request)
	if err != nil {
		logCtx(reqCtx).Warn().Err(err).
			Msg("Request context")
		resp.WriteHeaderAndJson(http.StatusInternalServerError, &rest.ErrorResponse{}, restful.MIME_JSON)
		return
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsUserSubject() {
		logCtx(reqCtx).Warn().Msg("Already authorized")
		resp.WriteHeaderAndJson(http.StatusOK,
			&authGetResponse{AccessToken: ctxAuth.RawToken()},
			restful.MIME_JSON)
		return
	}

	authCode := req.QueryParameter("code")

	accessToken, err := restSvc.iamCS.
		AccessTokenByAuthorizationCodeGrant(authCode)
	if err != nil {
		panic(err)
	}

	resp.WriteHeaderAndJson(http.StatusOK,
		&authGetResponse{AccessToken: accessToken},
		restful.MIME_JSON)
}

type helloGetResponse struct {
	Greetings string `json:"greetings"`
}

func (restSvc *RESTService) getHello(req *restful.Request, resp *restful.Response) {
	reqCtx, err := restSvc.iamCS.RESTCallInputContext(req.Request)
	if err != nil {
		logCtx(reqCtx).Warn().Err(err).
			Msg("Request context")
		resp.WriteHeaderAndJson(http.StatusInternalServerError, &rest.ErrorResponse{}, restful.MIME_JSON)
		return
	}
	ctxAuth := reqCtx.Authorization()
	if !ctxAuth.IsUserSubject() {
		logCtx(reqCtx).
			Warn().Msg("Unauthorized")
		resp.WriteHeaderAndJson(http.StatusUnauthorized, &rest.ErrorResponse{},
			restful.MIME_JSON)
		return
	}

	resp.WriteHeaderAndJson(http.StatusOK,
		&helloGetResponse{Greetings: "Hello, user " + ctxAuth.UserID().AZIDText()},
		restful.MIME_JSON)
}
