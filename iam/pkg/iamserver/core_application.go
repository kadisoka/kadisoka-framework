package iamserver

import (
	"crypto/subtle"

	errors "github.com/alloyzeus/go-azfl/errors"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) ApplicationByID(id iam.ApplicationID) (*iam.Application, error) {
	if core.applicationDataProvider == nil {
		return nil, nil
	}
	return core.applicationDataProvider.GetApplication(id)
}

// AuthenticatedApplication returns an instance of application only when the
// authentication succeeded.
func (core *Core) AuthenticatedApplication(
	appID iam.ApplicationID, secret string,
) (*iam.Application, error) {
	app, err := core.ApplicationByID(appID)
	if err != nil {
		return nil, errors.Wrap("app look up", err)
	}
	if app == nil {
		return nil, nil
	}
	if subtle.ConstantTimeCompare([]byte(secret), []byte(app.Data.Secret)) != 1 {
		return nil, errors.ArgMsg("password", "mismatch")
	}

	return app, nil
}
