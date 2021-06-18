package iam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	grpcmd "google.golang.org/grpc/metadata"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/api/oauth2"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
)

func NewClientAppSimple(
	appBase app.App,
	envVarsPrefix string,
) (*ServiceClientApp, error) {
	svc, err := NewServiceClientSimple(appBase.InstanceID(), envVarsPrefix)
	if err != nil {
		return nil, errors.Wrap("service client initialization", err)
	}

	return &ServiceClientApp{
		App:           appBase,
		ServiceClient: svc,
	}, nil
}

type ServiceClientApp struct {
	app.App
	ServiceClient
}

// ServiceClient is an abstraction to be used by applications to access
// IAM resources.
//
// This interface is not designed for service applications as it doesn't
// have the required abstractions to handle access from other applications.
// For interface which could be used by service applications, see
// ConsumerServer .
type ServiceClient interface {
	// ServerBaseURL returns the base URL of the IAM server this client
	// will connect to.
	ServerBaseURL() string

	// TerminalRef returns the terminal ref-key of the client instance after
	// successful authentication with IAM server.
	TerminalRef() TerminalRefKey

	GRPCServiceClient
	RESTServiceClient

	ServiceClientAuth

	UserServiceClient
}

type ServiceClientAuth interface {
	// AuthenticateServiceClient authenticates current application as a
	// service which will grant access to S2S API as configured on the
	// IAM service server.
	AuthenticateServiceClient(
		serviceInstanceID string,
	) (terminalRef TerminalRefKey, err error)

	// AccessTokenByAuthorizationCodeGrant obtains access token by providing
	// authorization code returned from a 3-legged authorization flow
	// (the authorization code flow).
	AccessTokenByAuthorizationCodeGrant(
		authorizationCode string,
	) (accessToken string, err error)
}

const (
	serverOAuth2JWKSRelPath  = "/oauth2/jwks"
	serverOAuth2TokenRelPath = "/oauth2/token"
)

func NewServiceClientSimple(
	instID string,
	envVarsPrefix string,
) (ServiceClient, error) {
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

	inst, err := NewServiceClient(cfg, &jwtKeyChain, userInstanceInfoService)
	if err != nil {
		return nil, err
	}

	_, err = inst.AuthenticateServiceClient(instID)
	if err != nil {
		return nil, err
	}

	return inst, nil
}

func NewServiceClient(
	serviceClientConfig *ServiceClientConfig,
	jwtKeyChain *JWTKeyChain,
	userInstanceInfoService UserInstanceInfoService,
) (ServiceClient, error) {
	if serviceClientConfig != nil {
		cfg := *serviceClientConfig
		serviceClientConfig = &cfg
	}

	return &serviceClientCore{
		serviceClientConfig: serviceClientConfig,
		userInstanceInfoSvc: userInstanceInfoService,
	}, nil
}

type serviceClientCore struct {
	serviceClientConfig *ServiceClientConfig
	terminalRef         TerminalRefKey
	clientAccessToken   string

	//HACK: implement in this struct rather than forwarding it
	userInstanceInfoSvc UserInstanceInfoService
}

var _ ServiceClient = &serviceClientCore{}

func (svcClient *serviceClientCore) ServerBaseURL() string {
	if svcClient.serviceClientConfig != nil {
		return svcClient.serviceClientConfig.ServerBaseURL
	}
	return ""
}

func (svcClient *serviceClientCore) TerminalRef() TerminalRefKey { return svcClient.terminalRef }

func (svcClient *serviceClientCore) AuthenticateServiceClient(
	serviceInstanceID string,
) (terminalRef TerminalRefKey, err error) {
	if svcClient.serviceClientConfig == nil {
		return TerminalRefKeyZero(), errors.New("oauth client is not configured")
	}
	baseURL := svcClient.ServerBaseURL()
	if !strings.HasPrefix(baseURL, "http") {
		return TerminalRefKeyZero(), errors.New("iam server base URL is not configured")
	}

	if serviceInstanceID == "" {
		return TerminalRefKeyZero(), errors.ArgMsg("serviceInstanceID", "empty")
	}

	terminalRef, accessToken, err := svcClient.
		obtainAccessTokenByClientCredentials(serviceInstanceID)
	if err != nil {
		panic(err)
	}

	svcClient.terminalRef = terminalRef
	svcClient.clientAccessToken = accessToken

	return svcClient.terminalRef, nil
}

func (svcClient *serviceClientCore) obtainAccessTokenByClientCredentials(
	serviceInstanceID string,
) (terminalRef TerminalRefKey, accessToken string, err error) {
	if svcClient.serviceClientConfig == nil || svcClient.serviceClientConfig.Credentials.ClientID == "" {
		return TerminalRefKeyZero(), "", errors.New("oauth client is not configured")
	}
	baseURL := svcClient.ServerBaseURL()
	if !strings.HasPrefix(baseURL, "http") {
		return TerminalRefKeyZero(), "", errors.New("iam server base URL is not configured")
	}
	tokenEndpointURL := baseURL + serverOAuth2TokenRelPath

	payloadStr, err := oauth2.QueryString(oauth2.AccessTokenRequest{
		GrantType: oauth2.GrantTypeClientCredentials,
	})
	if err != nil {
		return TerminalRefKeyZero(), "", errors.Wrap("outgoing request encoding", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		tokenEndpointURL,
		bytes.NewBuffer([]byte(payloadStr)))
	if err != nil {
		return TerminalRefKeyZero(), "", err
	}

	req.SetBasicAuth(
		svcClient.serviceClientConfig.Credentials.ClientID,
		svcClient.serviceClientConfig.Credentials.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	goRuntimeVersion := runtime.Version()
	goRuntimeVersion = "go/" + strings.TrimPrefix(goRuntimeVersion, "go")
	// include app.Info?
	req.Header.Set("User-Agent", "Kadisoka-IAM-Client/1.0 ("+serviceInstanceID+") "+goRuntimeVersion)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TerminalRefKeyZero(), "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Unexpected response status: " + resp.Status)
	}

	var tokenResp OAuth2TokenResponse
	err = json.NewDecoder(resp.Body).
		Decode(&tokenResp)
	if err != nil {
		return TerminalRefKeyZero(), "", err
	}

	terminalRef, err = TerminalRefKeyFromAZIDText(tokenResp.TerminalID)
	if err != nil {
		return TerminalRefKeyZero(), "", errors.Wrap("TerminalRefKeyFromAZIDText", err)
	}

	//TODO: to handle expiration, we'll need to store the value of 'ExpiresIn'
	// from the response[1] or 'exp' from the JWT claims[2].
	// [1] https://tools.ietf.org/html/rfc6749#section-4.2.2
	// [2] https://tools.ietf.org/html/rfc7519#section-4.1.4

	return terminalRef, tokenResp.AccessToken, nil
}

func (svcClient *serviceClientCore) getClientAccessToken() string {
	//TOOD:
	// - check the expiration. if the token is about to expire, 1 minute
	//   before expiration which info was obtained in obtainAccessTokenByPasswordWithTerminalCreds,
	//   start a task (goroutine) to obtain a new token
	// - ensure that only one task running at a time (mutex)
	return svcClient.clientAccessToken
}

// AccessTokenByAuthorizationCodeGrant conforms ServiceClientAuth.
func (svcClient *serviceClientCore) AccessTokenByAuthorizationCodeGrant(
	authorizationCode string,
) (accessToken string, err error) {
	if svcClient.serviceClientConfig == nil {
		return "", errors.New("oauth client is not configured")
	}
	baseURL := svcClient.ServerBaseURL()
	if !strings.HasPrefix(baseURL, "http") {
		return "", errors.New("iam server base URL is not configured")
	}
	tokenEndpointURL := baseURL + serverOAuth2TokenRelPath

	//TODO: redirect_uri is required
	payloadStr, err := oauth2.QueryString(oauth2.AccessTokenRequest{
		GrantType: oauth2.GrantTypeAuthorizationCode,
		Code:      authorizationCode,
	})
	if err != nil {
		return "", errors.Wrap("outgoing request encoding", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		tokenEndpointURL,
		bytes.NewBuffer([]byte(payloadStr)))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(
		svcClient.serviceClientConfig.Credentials.ClientID,
		svcClient.serviceClientConfig.Credentials.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp oauth2.ErrorResponse
		err = json.NewDecoder(resp.Body).
			Decode(&errResp)
		return "", errors.Msg(fmt.Sprintf("unable to exchange authorization code with access token: %v - %v",
			errResp.Error, errResp.ErrorDescription))
	}

	var tokenResp oauth2.TokenResponse
	err = json.NewDecoder(resp.Body).
		Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	//TODO: to handle expiration, we'll need to store the value of 'ExpiresIn'
	// from the response[1] or 'exp' from the JWT claims[2].
	// [1] https://tools.ietf.org/html/rfc6749#section-4.2.2
	// [2] https://tools.ietf.org/html/rfc7519#section-4.1.4

	return tokenResp.AccessToken, nil
}

// AuthorizedOutgoingGRPCContext returns a new instance of Context with
// authorization information set. If baseContext is valid, this method
// will use it as the parent context, otherwise, this method will create
// a Background context.
func (svcClient *serviceClientCore) AuthorizedOutgoingGRPCContext(
	baseContext context.Context,
) context.Context {
	accessToken := svcClient.getClientAccessToken()
	md := grpcmd.Pairs(AuthorizationMetadataKey, accessToken)
	if baseContext == nil {
		baseContext = context.Background()
	}
	return grpcmd.NewOutgoingContext(baseContext, md)
}

// AuthorizedOutgoingHTTPRequestHeader returns a new instance of http.Header
// with authorization information set. If baseHeader is proivded, this method
// will merge it into the returned value.
func (svcClient *serviceClientCore) AuthorizedOutgoingHTTPRequestHeader(
	baseHeader http.Header,
) http.Header {
	accessToken := svcClient.getClientAccessToken()
	outHeader := http.Header{}
	if accessToken != "" {
		outHeader.Set("Authorization", "Bearer "+accessToken)
	}
	if len(baseHeader) > 0 {
		for k, v := range baseHeader {
			outHeader[k] = v[:]
		}
	}
	return outHeader
}

func (svcClient *serviceClientCore) GetUserInstanceInfo(
	callCtx OpInputContext,
	userRef UserRefKey,
) (*UserInstanceInfo, error) {
	return svcClient.userInstanceInfoSvc.
		GetUserInstanceInfo(callCtx, userRef)
}
