package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
	restfulopenapi "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/rest"
	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/app"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam/rest/sec"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/rest/oauth2"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/rest/terminal"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/rest/user"
)

const ServerLatestVersionString = "v1"

type ServerConfig struct {
	ServePort int    `env:"SERVE_PORT"`
	ServePath string `env:"SERVE_PATH"`

	// SwaggerUIAssetsDir provides information where the Swagger UI files are
	// located. If left empty, the server won't serve the Swagger UI.
	SwaggerUIAssetsDir string `env:"SWAGGER_UI_ASSETS_DIR"`

	// V1 contains configuration for version 1 of the API service
	V1 *ServerV1Config `env:"V1"`
}

// Normalized returns a normalized copy
func (serverCfg ServerConfig) Normalized() ServerConfig {
	out := serverCfg
	out.ServePath = strings.TrimSuffix(out.ServePath, "/")
	return out
}

func (serverCfg ServerConfig) APISpecServePath() string {
	return serverCfg.ServePath + "/apidocs.json"
}

type ServerV1Config struct {
	ServePath string `env:"SERVE_PATH"`
}

type Server struct {
	config       ServerConfig
	httpServer   *http.Server
	serveMux     *http.ServeMux
	shuttingDown bool
}

var serviceInfo = app.ServiceInfo{
	Name:        "IAM REST API",
	Description: "Identity and Access Management service REST API",
}

// ServiceInfo conforms app.ServiceServer interface.
func (srv *Server) ServiceInfo() app.ServiceInfo { return serviceInfo }

// Serve conforms app.ServiceServer interface.
func (srv *Server) Serve() error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.config.ServePort),
		Handler: srv.serveMux}
	srv.httpServer = httpServer
	err := srv.httpServer.ListenAndServe()
	if err == nil {
		if !srv.shuttingDown {
			return errors.Msg("server stopped unexpectedly")
		}
		return nil
	}
	if err == http.ErrServerClosed && srv.shuttingDown {
		return nil
	}
	return err
}

// Shutdown conforms app.ServiceServer interface.
func (srv *Server) Shutdown(ctx context.Context) error {
	//TODO: mutex?
	srv.shuttingDown = true
	return srv.httpServer.Shutdown(ctx)
}

// IsAcceptingClients conforms app.ServiceServer interface.
func (srv Server) IsAcceptingClients() bool {
	return !srv.shuttingDown && srv.IsHealthy()
}

// IsHealthy conforms app.ServiceServer interface.
func (srv Server) IsHealthy() bool { return true }

// ServeHTTP conforms Go's HTTP Handler interface.
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.serveMux.ServeHTTP(w, r)
}

func (srv *Server) handleHealthCheck(w http.ResponseWriter, _ *http.Request) {
	if !srv.IsHealthy() {
		log.Error().Msg("Service is not healthy")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}

func NewServer(
	appApp app.App,
	config ServerConfig,
	iamServerCore *iamserver.Core,
	webUIURLs *iam.WebUIURLs, //TODO: add this to config (not sure if it's server core or ServerConfig)
) (*Server, error) {
	if webUIURLs == nil || webUIURLs.SignIn == "" {
		return nil, errors.Msg("requires valid web UI config")
	}

	config = config.Normalized()

	servePath := config.ServePath
	apiSpecServePath := config.APISpecServePath()

	serveMux := http.NewServeMux()
	restfulContainer, err := setUpRestfulContainer(config, serveMux, iamServerCore, servePath, webUIURLs)
	if err != nil {
		return nil, errors.Wrap("REST API container set up", err)
	}

	restfulSpecConfig, err := setUpRestfulSpecConfig(appApp.AppInfo(), apiSpecServePath, restfulContainer)
	if err != nil {
		return nil, errors.Wrap("REST API OpenAPI spec config set up", err)
	}

	// Setup API specification handler
	restfulContainer.Add(restfulopenapi.NewOpenAPIService(*restfulSpecConfig))

	log.Info().Msgf("REST API spec at %s", apiSpecServePath)
	if config.SwaggerUIAssetsDir != "" {
		// The trailing slash here is important.
		apiDocsUIPath := servePath + "/apidocs/"
		serveMux.Handle(apiDocsUIPath,
			http.StripPrefix(apiDocsUIPath,
				http.FileServer(http.Dir(config.SwaggerUIAssetsDir))))
		log.Info().Msgf("REST API documentations UI at %s", apiDocsUIPath)
	}

	srv := &Server{
		config:   config,
		serveMux: serveMux,
	}

	// Health check is used by load balancer and/or orchestrator
	serveMux.HandleFunc(
		"/healthz", srv.handleHealthCheck)

	return srv, nil
}

func GenerateOpenAPISpec(
	appApp app.App,
	config ServerConfig,
	iamServerCore *iamserver.Core,
	webUIURLs *iam.WebUIURLs, //TODO: add this to config (not sure if it's server core or ServerConfig)
) (*spec.Swagger, error) {
	config = config.Normalized()

	servePath := config.ServePath
	apiSpecServePath := config.APISpecServePath()

	restfulContainer, err := setUpRestfulContainer(config, nil, iamServerCore, servePath, webUIURLs)
	if err != nil {
		return nil, errors.Wrap("REST API container set up", err)
	}

	restfulSpecConfig, err := setUpRestfulSpecConfig(appApp.AppInfo(), apiSpecServePath, restfulContainer)
	if err != nil {
		return nil, errors.Wrap("REST API OpenAPI spec config set up", err)
	}

	spec := restfulopenapi.BuildSwagger(*restfulSpecConfig)

	return spec, nil
}

func initRESTV1Services(
	servePath string,
	container *restful.Container,
	iamServerCore *iamserver.Core,
	signInURL string,
) {
	log.Info().Msg("Initializing terminal service...")
	terminalSrv := terminal.NewServer(iamServerCore,
		terminal.ServerConfig{
			ServePath: servePath + "/terminals",
		})
	container.Add(terminalSrv.RestfulWebService())

	log.Info().Msg("Initializing user service...")
	userSrv := user.NewServer(iamServerCore,
		user.ServerConfig{
			ServePath: servePath + "/users",
		})
	container.Add(userSrv.RestfulWebService())

	log.Info().Msg("Initializing OAuth 2.0 service...")
	oauth2Srv, err := oauth2.NewServer(iamServerCore,
		oauth2.ServerConfig{
			ServePath: servePath + "/oauth2",
			SignInURL: signInURL,
		})
	if err != nil {
		log.Fatal().Err(err).
			Msg("OAuth 2.0 service initialization")
	}
	container.Add(oauth2Srv.RestfulWebService())
}

func processSwaggerSpec(
	buildInfo app.BuildInfo,
	swaggerSpec *spec.Swagger,
	secDefs spec.SecurityDefinitions,
) {
	rev := buildInfo.RevisionID
	if rev != "unknown" && len(rev) > 7 {
		rev = rev[:7]
	}
	swaggerSpec.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       serviceInfo.Name,
			Description: serviceInfo.Description,
			Version: fmt.Sprintf(
				"0.0.0-%s built at %s",
				rev, buildInfo.Timestamp),
		},
	}
	for k := range swaggerSpec.Paths.Paths {
		swaggerSpec.Paths.Paths[k] = processOpenAPIPath(
			swaggerSpec.Paths.Paths[k], secDefs)
	}
	swaggerSpec.SecurityDefinitions = secDefs
}

func setUpRestfulSpecConfig(
	appInfo app.Info,
	apiSpecServePath string,
	restfulContainer *restful.Container,
) (*restfulopenapi.Config, error) {
	return &restfulopenapi.Config{
		WebServices: restfulContainer.RegisteredWebServices(),
		APIPath:     apiSpecServePath,
		PostBuildSwaggerObjectHandler: func(swaggerSpec *spec.Swagger) {
			processSwaggerSpec(appInfo.BuildInfo, swaggerSpec, securityDefinitionsDefault())
		},
	}, nil
}

func securityDefinitionsDefault() spec.SecurityDefinitions {
	return _securityDefinitionsDefault
}

var _securityDefinitionsDefault = spec.SecurityDefinitions{
	sec.AuthorizationBasicOAuth2ClientCredentials.String(): spec.BasicAuth(),
	//TODO: change to spec.OAuth2AccessToken (?)
	sec.AuthorizationBearerAccessToken.String(): spec.APIKeyAuth(iam.AuthorizationMetadataKey, "header"),
}

func setUpRestfulContainer(
	config ServerConfig,
	serveMux *http.ServeMux,
	iamServerCore *iamserver.Core,
	servePath string,
	webUIURLs *iam.WebUIURLs,
) (*restful.Container, error) {
	container := restful.NewContainer()
	container.ServeMux = serveMux
	container.EnableContentEncoding(true)
	container.ServiceErrorHandler(
		func(
			err restful.ServiceError,
			req *restful.Request,
			resp *restful.Response,
		) {
			logReq(req.Request).
				Warn().Int("status_code", err.Code).Str("err_msg", err.Message).
				Msg("Routing error")
			resp.WriteErrorString(err.Code, err.Message)
		})

	// We need CORS for our webclients
	rest.SetUpCORSFilterByEnv(container, "CORS_", nil) //TODO: from config

	var v1ServePath string
	if config.V1 != nil {
		v1ServePath = config.V1.ServePath
	}
	if v1ServePath == "" {
		v1ServePath = servePath + "/v1"
	}

	initRESTV1Services(v1ServePath, container, iamServerCore, webUIURLs.SignIn)

	return container, nil
}

func processOpenAPIPath(
	pathItem spec.PathItem, secDefs spec.SecurityDefinitions,
) spec.PathItem {
	pathItem.Get = processOpenAPIPathOp(pathItem.Get, secDefs)
	pathItem.Put = processOpenAPIPathOp(pathItem.Put, secDefs)
	pathItem.Post = processOpenAPIPathOp(pathItem.Post, secDefs)
	pathItem.Delete = processOpenAPIPathOp(pathItem.Delete, secDefs)
	pathItem.Options = processOpenAPIPathOp(pathItem.Options, secDefs)
	pathItem.Head = processOpenAPIPathOp(pathItem.Head, secDefs)
	pathItem.Patch = processOpenAPIPathOp(pathItem.Patch, secDefs)
	return pathItem
}

func processOpenAPIPathOp(
	op *spec.Operation, secDefs spec.SecurityDefinitions,
) *spec.Operation {
	if op == nil {
		return nil
	}

	for _, tag := range op.Tags {
		if tag == "hidden" {
			return nil
		}
	}

	var updatedParams []spec.Parameter
	for _, p := range op.Parameters {
		isSec := false
		if p.Description != "" {
			lowerDesc := strings.ToLower(p.Description)
			for k, secDef := range secDefs {
				if strings.HasPrefix(lowerDesc, k) {
					if secDef.Type == "basic" {
						// Basic authorization is always as 'Authorization' in the header
						if p.Name == iam.AuthorizationMetadataKey && p.In == "header" {
							op.Security = append(op.Security, map[string][]string{k: {}})
							isSec = true
							continue
						}
					}
					if secDef.Type == "apiKey" || secDef.Type == "accessCode" {
						if p.Name == secDef.Name && p.In == secDef.In {
							op.Security = append(op.Security, map[string][]string{k: {}})
							isSec = true
							continue
						}
					}
				}
			}
		}
		if !isSec {
			updatedParams = append(updatedParams, p)
		}
	}
	op.Parameters = updatedParams
	return op
}
