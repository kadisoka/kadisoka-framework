package iamserver

import (
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) ApplicationByRefKey(refKey iam.ApplicationRefKey) (*iam.Client, error) {
	if core.clientDataProvider == nil {
		return nil, nil
	}
	return core.clientDataProvider.GetClient(refKey)
}
