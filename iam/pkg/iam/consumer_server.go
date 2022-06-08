package iam

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/alloyzeus/go-azfl/azcore"
	"github.com/alloyzeus/go-azfl/errors"
	dataerrs "github.com/alloyzeus/go-azfl/errors/data"
	"github.com/square/go-jose/v3/jwt"
	"github.com/tomasen/realip"
	"golang.org/x/text/language"
	grpcmd "google.golang.org/grpc/metadata"
	grpcpeer "google.golang.org/grpc/peer"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
)

func NewConsumerServerAppSimple(
	appBase app.App,
	envVarsPrefix string,
) (*ConsumerServerApp, error) {
	svc, err := NewConsumerServerSimple(appBase.InstanceID(), envVarsPrefix)
	if err != nil {
		return nil, errors.Wrap("service client initialization", err)
	}

	return &ConsumerServerApp{
		App:            appBase,
		ConsumerServer: svc,
	}, nil
}

type ConsumerServerApp struct {
	app.App
	ConsumerServer
}

func NewConsumerServerSimple(
	instID string,
	envVarsPrefix string,
) (ConsumerServer, error) {
	cfg, err := ServiceClientConfigFromEnv(envVarsPrefix, nil)
	if err != nil {
		return nil, errors.Wrap("config loading", err)
	}

	jwksURL := cfg.ServerBaseURL + serverOAuth2JWKSRelPath
	var jwtKeyChain JWTKeyChain
	_, err = jwtKeyChain.LoadVerifierKeysFromJWKSetByURL(jwksURL)
	if err != nil {
		return nil, errors.Wrap("jwt key set loading", err)
	}

	userInstanceInfoService := &UserInstanceInfoServiceClientCore{}

	inst, err := NewConsumerServer(cfg, &jwtKeyChain, userInstanceInfoService)
	if err != nil {
		return nil, err
	}

	_, err = inst.AuthenticateServiceClient(instID)
	if err != nil {
		return nil, err
	}

	return inst, nil
}

func NewConsumerServer(
	serviceClientConfig *ServiceClientConfig,
	jwtKeyChain *JWTKeyChain,
	userInstanceInfoService UserInstanceInfoService,
) (ConsumerServer, error) {
	if serviceClientConfig != nil {
		cfg := *serviceClientConfig
		serviceClientConfig = &cfg
	}

	srvCore, err := NewConsumerServerCore(jwtKeyChain, userInstanceInfoService)
	if err != nil {
		return nil, err
	}

	return &consumerServerCore{
		&serviceClientCore{
			serviceClientConfig: serviceClientConfig,
			userInstanceInfoSvc: userInstanceInfoService,
		},
		srvCore,
	}, nil
}

// ConsumerServer is an abstractions for a server which acts as
// a client/consumer of IAM, and also allow applications authorized by IAM
// to access its resources.
type ConsumerServer interface {
	ConsumerServerBase
	ServiceClient
}

// ConsumerServerBase is an interface which contains utilities for
// IAM service clients to handle requests from other IAM service clients.
type ConsumerServerBase interface {
	// AuthorizationFromJWTString loads authorization context from a JWT
	// string.
	AuthorizationFromJWTString(
		jwtStr string,
	) (*Authorization, error)

	// JWTKeyChain returns instance of key chain used to sign JWT tokens.
	JWTKeyChain() *JWTKeyChain

	ConsumerGRPCServer
	ConsumerRESTServer
}

// ConsumerGRPCServer is an interface which contains utilities for
// IAM service clients to handle requests from other clients.
type ConsumerGRPCServer interface {
	// GRPCOpInputContext loads authorization context from
	// gRPC call context.
	GRPCOpInputContext(
		grpcContext context.Context,
	) (*GRPCOpInputContext, error)
}

// ConsumerRESTServer is an interface which contains utilities for
// IAM service clients to handle requests from other clients.
type ConsumerRESTServer interface {
	// RESTOpInputContext returns a RESTOpInputContext instance for the request.
	// This function will always return an instance even if there's an error.
	RESTOpInputContext(*http.Request) (*RESTOpInputContext, error)
}

type consumerServerCore struct {
	*serviceClientCore
	ConsumerServerBase
}

func NewConsumerServerCore(
	jwtKeyChain *JWTKeyChain,
	userInstanceInfoService UserInstanceInfoService,
) (ConsumerServerBase, error) {
	return &consumerServerBaseCore{
		jwtKeyChain:             jwtKeyChain,
		userInstanceInfoService: userInstanceInfoService,
	}, nil
}

type consumerServerBaseCore struct {
	jwtKeyChain             *JWTKeyChain
	userInstanceInfoService UserInstanceInfoService
}

var _ ConsumerServerBase = &consumerServerBaseCore{}

func (consumerSrv *consumerServerBaseCore) JWTKeyChain() *JWTKeyChain {
	return consumerSrv.jwtKeyChain
}

// Shortcut
func (consumerSrv *consumerServerBaseCore) GetSignedVerifierKey(keyID string) interface{} {
	return consumerSrv.jwtKeyChain.GetSignedVerifierKey(keyID)
}

func (consumerSrv *consumerServerBaseCore) AuthorizationFromJWTString(
	jwtStr string,
) (*Authorization, error) {
	emptyAuthCtx := newEmptyAuthorization()
	if jwtStr == "" {
		return emptyAuthCtx, nil
	}

	tok, err := jwt.ParseSigned(jwtStr)
	if err != nil {
		return emptyAuthCtx, errors.ArgWrap("", "parsing", err)
	}
	if len(tok.Headers) != 1 {
		return emptyAuthCtx, errors.ArgMsg("", "invalid number of headers")
	}

	keyID := tok.Headers[0].KeyID
	if keyID == "" {
		return emptyAuthCtx, errors.Arg("", errors.EntMsg("kid", "empty"))
	}

	verifierKey := consumerSrv.JWTKeyChain().GetSignedVerifierKey(keyID)
	if verifierKey == nil {
		return emptyAuthCtx, errors.Arg("", errors.EntMsg("kid", "reference invalid"))
	}

	var claims AccessTokenClaims
	err = tok.Claims(verifierKey, &claims)
	if err != nil {
		return emptyAuthCtx, errors.ArgWrap("", "verification", err)
	}

	//TODO: check expiry

	if claims.ID == "" {
		return emptyAuthCtx, errors.Arg("", errors.EntMsg("jti", "empty"))
	}
	sessionRef, err := SessionRefKeyFromAZIDText(claims.ID)
	if err != nil {
		return emptyAuthCtx, errors.Arg("", errors.Ent("jti", dataerrs.Malformed(err)))
	}
	//TODO(exa): check if the authorization instance id has been revoked

	var userRef UserRefKey
	if claims.Subject != "" {
		userRef, err = UserRefKeyFromAZIDText(claims.Subject)
		if err != nil {
			return emptyAuthCtx, errors.Arg("", errors.EntMsg("sub", "malformed"))
		}
		instInfo, err := consumerSrv.userInstanceInfoService.
			GetUserInstanceInfo(nil, userRef)
		if err != nil {
			return emptyAuthCtx, errors.Wrap("account state query", err)
		}
		if instInfo == nil {
			return emptyAuthCtx, errors.Arg("", errors.EntMsg("sub", "reference invalid"))
		}
		if !instInfo.IsActive() {
			return emptyAuthCtx, errors.Arg("", errors.EntMsg("sub", "reference invalid"))
		}
	}

	var terminalRef TerminalRefKey
	if claims.TerminalID == "" {
		return emptyAuthCtx, errors.Arg("", errors.EntMsg("terminal_id", "empty"))
	}
	terminalRef, err = TerminalRefKeyFromAZIDText(claims.TerminalID)
	if err != nil {
		return emptyAuthCtx, errors.Arg("", errors.Ent("terminal_id", dataerrs.Malformed(err)))
	}
	if terminalRef.IsNotStaticallyValid() {
		return emptyAuthCtx, errors.Arg("", errors.Ent("terminal_id", dataerrs.ErrMalformed))
	}

	return &Authorization{
		Session:  sessionRef,
		rawToken: jwtStr,
	}, nil
}

func (consumerSrv *consumerServerBaseCore) GRPCOpInputContext(
	grpcCallCtx context.Context,
) (*GRPCOpInputContext, error) {
	callCtx, err := consumerSrv.callContextFromGRPCContext(grpcCallCtx)
	if callCtx == nil {
		callCtx = NewEmptyOpInputContext(grpcCallCtx)
	}
	return &GRPCOpInputContext{callCtx}, err
}

func (consumerSrv *consumerServerBaseCore) callContextFromGRPCContext(
	grpcCallCtx context.Context,
) (OpInputContext, error) {
	var remoteAddr string
	if peer, _ := grpcpeer.FromContext(grpcCallCtx); peer != nil {
		remoteAddr = peer.Addr.String() //TODO: attempt to resolve if it's proxied
	}

	var originEnvString string
	var originAcceptLanguages []language.Tag

	if md, mdOK := grpcmd.FromIncomingContext(grpcCallCtx); mdOK {
		userAgentMDVal := md.Get("user-agent")
		if len(userAgentMDVal) > 0 {
			originEnvString = userAgentMDVal[0]
		}

		acceptLanguageMDVal := md.Get("accept-language")
		if len(acceptLanguageMDVal) > 0 {
			originAcceptLanguages, _, _ = language.ParseAcceptLanguage(acceptLanguageMDVal[0])
		}
	}

	originInfo := azcore.ServiceMethodCallOriginInfo{
		Address:           remoteAddr,
		AcceptLanguage:    originAcceptLanguages,
		EnvironmentString: originEnvString,
	}

	ctxAuth, err := consumerSrv.authorizationFromGRPCContext(grpcCallCtx)
	if err != nil {
		return newOpInputContext(grpcCallCtx, ctxAuth, originInfo, nil), err
	}

	var opID *api.OpID
	md, ok := grpcmd.FromIncomingContext(grpcCallCtx)
	if !ok {
		return newOpInputContext(grpcCallCtx, ctxAuth, originInfo, nil), nil
	}

	//TODO: idempotency key https://datatracker.ietf.org/doc/html/draft-ietf-httpapi-idempotency-key-header-01
	opIDStrs := md.Get("op-id")
	if len(opIDStrs) == 0 {
		opIDStrs = md.Get("request-id")
		if len(opIDStrs) == 0 {
			opIDStrs = md.Get("x-request-id")
		}
	}
	if len(opIDStrs) > 0 {
		opIDStr := opIDStrs[0]
		i, err := api.OpIDFromString(opIDStr)
		if err != nil {
			return newOpInputContext(grpcCallCtx, ctxAuth, originInfo, nil),
				ReqFieldErr("Request-ID", err)
		}
		opID = &i
	}

	return newOpInputContext(grpcCallCtx, ctxAuth, originInfo, opID), err
}

func (consumerSrv *consumerServerBaseCore) authorizationFromGRPCContext(
	grpcContext context.Context,
) (*Authorization, error) {
	emptyAuthCtx := newEmptyAuthorization()
	md, ok := grpcmd.FromIncomingContext(grpcContext)
	if !ok {
		return emptyAuthCtx, nil
	}
	authorizations := md.Get(AuthorizationMetadataKey)
	if len(authorizations) == 0 {
		return emptyAuthCtx, nil
	}
	if authorizations[0] == "" {
		return emptyAuthCtx, ReqFieldErr("Authorization", dataerrs.ErrEmpty)
	}
	token := authorizations[0]
	parts := strings.SplitN(token, " ", 2)
	if len(parts) == 2 {
		if strings.ToLower(parts[0]) != "bearer" {
			return emptyAuthCtx, ErrReqFieldAuthorizationTypeUnsupported
		}
		token = parts[1]
	}
	return consumerSrv.AuthorizationFromJWTString(token)
}

// RESTOpInputContext creates a call context which represents an HTTP request.
func (consumerSrv *consumerServerBaseCore) RESTOpInputContext(
	req *http.Request,
) (*RESTOpInputContext, error) {
	callCtx, err := consumerSrv.callContextFromHTTPRequest(req)
	if callCtx == nil {
		callCtx = NewEmptyOpInputContext(req.Context())
	}
	return &RESTOpInputContext{callCtx, req}, err
}

func (consumerSrv *consumerServerBaseCore) callContextFromHTTPRequest(
	req *http.Request,
) (OpInputContext, error) {
	ctx := req.Context()
	ctxAuth := newEmptyAuthorization()

	remoteAddr := realip.FromRequest(req)
	if remoteAddr == "" {
		remoteAddr = req.RemoteAddr
	}

	remoteEnvString := req.UserAgent()
	acceptLanguages, _, _ := language.ParseAcceptLanguage(req.Header.Get("Accept-Language"))

	var originDateTime *time.Time
	if s := req.Header.Get("Date"); s != "" {
		dt, err := time.Parse(time.RFC1123, s)
		if err == nil {
			originDateTime = &dt
		}
	}

	originInfo := azcore.ServiceMethodCallOriginInfo{
		Address:           remoteAddr,
		EnvironmentString: remoteEnvString,
		AcceptLanguage:    acceptLanguages,
		DateTime:          originDateTime,
	}

	// Get from query too?
	var opID *api.OpID
	opIDStr := req.Header.Get("Op-ID")
	if opIDStr == "" {
		opIDStr = req.Header.Get("Request-ID")
		if opIDStr == "" {
			opIDStr = req.Header.Get("X-Request-ID")
		}
	}
	if opIDStr != "" {
		i, err := api.OpIDFromString(opIDStr)
		if err != nil {
			return newOpInputContext(ctx, ctxAuth, originInfo, nil),
				ReqFieldErr("Request-ID", err)
		}
		opID = &i
	}

	authorization := strings.TrimSpace(req.Header.Get("Authorization"))
	if authorization != "" {
		authParts := strings.SplitN(authorization, " ", 2)
		if len(authParts) != 2 {
			return newOpInputContext(ctx, ctxAuth, originInfo, nil),
				ErrReqFieldAuthorizationMalformed
		}
		if authParts[0] != "Bearer" {
			return newOpInputContext(ctx, ctxAuth, originInfo, nil),
				ErrReqFieldAuthorizationTypeUnsupported
		}

		jwtStr := strings.TrimSpace(authParts[1])
		var err error
		ctxAuth, err = consumerSrv.AuthorizationFromJWTString(jwtStr)
		if err != nil {
			return newOpInputContext(ctx, ctxAuth, originInfo, nil),
				ErrReqFieldAuthorizationMalformed
		}

		//TODO: validate ctxAuth
	}
	if ctxAuth == nil {
		ctxAuth = newEmptyAuthorization()
	}

	return newOpInputContext(ctx, ctxAuth, originInfo, opID), nil
}
