package iam

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
)

func NewAppSimple(envVarPrefix string) (*App, error) {
	appApp, err := app.InitByEnvDefault()
	if err != nil {
		return nil, errors.Wrap("app initialization", err)
	}

	iamClient, err := NewServiceClientSimple(appApp.InstanceID(), envVarPrefix)
	if err != nil {
		return nil, errors.Wrap("service client initialization", err)
	}

	return &App{
		App:           appApp,
		ServiceClient: iamClient,
	}, nil
}

type App struct {
	app.App
	ServiceClient
}
