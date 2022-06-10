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
	// GRPCCallInputContext loads authorization context from
	// gRPC call context.
	GRPCCallInputContext(
		grpcContext context.Context,
	) (*GRPCCallInputContext, error)
}

// ConsumerRESTServer is an interface which contains utilities for
// IAM service clients to handle requests from other clients.
type ConsumerRESTServer interface {
	// RESTCallInputContext returns a RESTCallInputContext instance for the request.
	// This function will always return an instance even if there's an error.
	RESTCallInputContext(*http.Request) (*RESTCallInputContext, error)
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
	sessionID, err := SessionIDFromAZIDText(claims.ID)
	if err != nil {
		return emptyAuthCtx, errors.Arg("", errors.Ent("jti", dataerrs.Malformed(err)))
	}
	//TODO(exa): check if the authorization instance id has been revoked

	var userID UserID
	if claims.Subject != "" {
		userID, err = UserIDFromAZIDText(claims.Subject)
		if err != nil {
			return emptyAuthCtx, errors.Arg("", errors.EntMsg("sub", "malformed"))
		}
		instInfo, err := consumerSrv.userInstanceInfoService.
			GetUserInstanceInfo(nil, userID)
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

	var terminalID TerminalID
	if claims.TerminalID == "" {
		return emptyAuthCtx, errors.Arg("", errors.EntMsg("terminal_id", "empty"))
	}
	terminalID, err = TerminalIDFromAZIDText(claims.TerminalID)
	if err != nil {
		return emptyAuthCtx, errors.Arg("", errors.Ent("terminal_id", dataerrs.Malformed(err)))
	}
	if terminalID.IsNotStaticallyValid() {
		return emptyAuthCtx, errors.Arg("", errors.Ent("terminal_id", dataerrs.ErrMalformed))
	}

	return &Authorization{
		Session:  sessionID,
		rawToken: jwtStr,
	}, nil
}

func (consumerSrv *consumerServerBaseCore) GRPCCallInputContext(
	grpcCallCtx context.Context,
) (*GRPCCallInputContext, error) {
	callCtx, err := consumerSrv.callContextFromGRPCContext(grpcCallCtx)
	if callCtx == nil {
		callCtx = NewEmptyCallInputContext(grpcCallCtx)
	}
	return &GRPCCallInputContext{callCtx}, err
}

func (consumerSrv *consumerServerBaseCore) callContextFromGRPCContext(
	grpcCallCtx context.Context,
) (CallInputContext, error) {
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
		return newCallInputContext(grpcCallCtx, ctxAuth, originInfo, nil), err
	}

	var idempotencyKey *api.IdempotencyKey
	md, ok := grpcmd.FromIncomingContext(grpcCallCtx)
	if !ok {
		return newCallInputContext(grpcCallCtx, ctxAuth, originInfo, nil), nil
	}

	//TODO: idempotency key https://datatracker.ietf.org/doc/html/draft-ietf-httpapi-idempotency-key-header-01
	idempotencyKeyStrs := md.Get("idempotency-key")
	if len(idempotencyKeyStrs) == 0 {
		idempotencyKeyStrs = md.Get("op-id")
	}
	if len(idempotencyKeyStrs) > 0 {
		keyStr := idempotencyKeyStrs[0]
		i, err := api.IdempotencyKeyFromString(keyStr)
		if err != nil {
			return newCallInputContext(grpcCallCtx, ctxAuth, originInfo, nil),
				ReqFieldErr("Idempotency-Key", err)
		}
		idempotencyKey = &i
	}

	return newCallInputContext(grpcCallCtx, ctxAuth, originInfo, idempotencyKey), err
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

// RESTCallInputContext creates a call context which represents an HTTP request.
func (consumerSrv *consumerServerBaseCore) RESTCallInputContext(
	req *http.Request,
) (*RESTCallInputContext, error) {
	callCtx, err := consumerSrv.callContextFromHTTPRequest(req)
	if callCtx == nil {
		callCtx = NewEmptyCallInputContext(req.Context())
	}
	return &RESTCallInputContext{callCtx, req}, err
}

func (consumerSrv *consumerServerBaseCore) callContextFromHTTPRequest(
	req *http.Request,
) (CallInputContext, error) {
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
	var idempotencyKey *api.IdempotencyKey
	idempotencyKeyStr := req.Header.Get("Idempotency-Key")
	if idempotencyKeyStr == "" {
		idempotencyKeyStr = req.Header.Get("Op-ID")
	}
	if idempotencyKeyStr != "" {
		i, err := api.IdempotencyKeyFromString(idempotencyKeyStr)
		if err != nil {
			return newCallInputContext(ctx, ctxAuth, originInfo, nil),
				ReqFieldErr("Idempotency-Key", err)
		}
		idempotencyKey = &i
	}

	authorization := strings.TrimSpace(req.Header.Get("Authorization"))
	if authorization != "" {
		authParts := strings.SplitN(authorization, " ", 2)
		if len(authParts) != 2 {
			return newCallInputContext(ctx, ctxAuth, originInfo, nil),
				ErrReqFieldAuthorizationMalformed
		}
		if authParts[0] != "Bearer" {
			return newCallInputContext(ctx, ctxAuth, originInfo, nil),
				ErrReqFieldAuthorizationTypeUnsupported
		}

		jwtStr := strings.TrimSpace(authParts[1])
		var err error
		ctxAuth, err = consumerSrv.AuthorizationFromJWTString(jwtStr)
		if err != nil {
			return newCallInputContext(ctx, ctxAuth, originInfo, nil),
				ErrReqFieldAuthorizationMalformed
		}

		//TODO: validate ctxAuth
	}
	if ctxAuth == nil {
		ctxAuth = newEmptyAuthorization()
	}

	return newCallInputContext(ctx, ctxAuth, originInfo, idempotencyKey), nil
}
