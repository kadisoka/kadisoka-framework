package iamserver

import "github.com/kadisoka/kadisoka-framework/iam/pkg/iam"

// Interface conformance assertion.
var _ iam.UserRefKeyService = &Core{}

// IsUserRefKeyRegistered is used to determine that a user ID has been registered.
// It's not checking if the account is active or not.
//
// This function is generally cheap if the user ID has been registered.
func (core *Core) IsUserRefKeyRegistered(refKey iam.UserRefKey) bool {
	idNum := refKey.IDNum()

	// Look up for an user ID in the cache.
	if _, idRegistered := core.registeredUserInstanceIDCache.Get(idNum); idRegistered {
		return true
	}

	idRegistered, _, err := core.
		getUserInstanceStateByIDNum(idNum)
	if err != nil {
		panic(err)
	}

	if idRegistered {
		core.registeredUserInstanceIDCache.Add(idNum, nil)
	}

	return idRegistered
}
