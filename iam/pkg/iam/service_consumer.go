package iam

import "github.com/alloyzeus/go-azfl/azfl/errors"

func NewServiceConsumerServerSimple(
	instID string,
	envVarsPrefix string,
) (ServiceConsumerServer, error) {
	cfg, err := ServiceClientConfigFromEnv(envVarsPrefix, nil)
	if err != nil {
		return nil, errors.Wrap("config loading", err)
	}

	jwksURL := cfg.ServerBaseURL + serverOAuth2JWKSPath
	var jwtKeyChain JWTKeyChain
	_, err = jwtKeyChain.LoadVerifierKeysFromJWKSetByURL(jwksURL)
	if err != nil {
		return nil, errors.Wrap("jwt key set loading", err)
	}

	uaStateServiceClient := &UserInstanceInfoServiceClientCore{}

	inst, err := NewServiceConsumerServer(cfg, &jwtKeyChain, uaStateServiceClient)
	if err != nil {
		return nil, err
	}

	_, err = inst.AuthenticateServiceClient(instID)
	if err != nil {
		return nil, err
	}

	return inst, nil
}

func NewServiceConsumerServer(
	serviceClientConfig *ServiceClientConfig,
	jwtKeyChain *JWTKeyChain,
	userInstanceInfoService UserInstanceInfoService,
) (ServiceConsumerServer, error) {
	if serviceClientConfig != nil {
		cfg := *serviceClientConfig
		serviceClientConfig = &cfg
	}

	serviceClientServer, err := NewServiceClientServer(jwtKeyChain, userInstanceInfoService)
	if err != nil {
		return nil, err
	}

	return &ServiceConsumerServerCore{
		&ServiceClientCore{
			serviceClientConfig: serviceClientConfig,
			userInstanceInfoSvc: userInstanceInfoService,
		},
		serviceClientServer,
	}, nil
}

type ServiceConsumerServerCore struct {
	*ServiceClientCore
	ServiceClientServer
}

// ServiceConsumerServer is an abstractions for a server which acts as
// an IAM client/consumer, and also allow applications, authorized by IAM
// to access its resources.
type ServiceConsumerServer interface {
	ServiceClientServer
	ServiceClient
}
