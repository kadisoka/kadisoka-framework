package iamserver

import (
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) ApplicationByID(id iam.ApplicationID) (*iam.Application, error) {
	if core.applicationDataProvider == nil {
		return nil, nil
	}
	return core.applicationDataProvider.GetApplication(id)
}
