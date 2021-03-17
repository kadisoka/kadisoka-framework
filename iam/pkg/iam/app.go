package iam

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
)

func NewAppSimple(envVarsPrefix string) (*App, error) {
	appApp := app.Instance()

	svc, err := NewServiceClientSimple(appApp.InstanceID(), envVarsPrefix)
	if err != nil {
		return nil, errors.Wrap("service client initialization", err)
	}

	return &App{
		App:           appApp,
		ServiceClient: svc,
	}, nil
}

func NewConsumerServerAppSimple(envVarsPrefix string) (*ConsumerServerApp, error) {
	appApp := app.Instance()

	svc, err := NewServiceConsumerServerSimple(appApp.InstanceID(), envVarsPrefix)
	if err != nil {
		return nil, errors.Wrap("service client initialization", err)
	}

	return &ConsumerServerApp{
		App:                   appApp,
		ServiceConsumerServer: svc,
	}, nil
}

type App struct {
	app.App
	ServiceClient
}

type ConsumerServerApp struct {
	app.App
	ServiceConsumerServer
}
