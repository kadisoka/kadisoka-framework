package iam

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
)

func NewAppSimple(envVarsPrefix string) (*App, error) {
	appApp := app.Instance()

	iamClient, err := NewServiceClientSimple(appApp.InstanceID(), envVarsPrefix)
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
