package iamserver

import "github.com/kadisoka/kadisoka-framework/iam/pkg/iam"

// Interface conformance assertion.
var _ iam.UserRefKeyService = &Core{}

// IsUserRefKeyRegistered is used to determine that a user ID has been registered.
// It's not checking if the account is active or not.
//
// This function is generally cheap if the user ID has been registered.
func (core *Core) IsUserRefKeyRegistered(refKey iam.UserRefKey) bool {
	instID := refKey.ID()

	// Look up for an user ID in the cache.
	if _, idRegistered := core.registeredUserIDCache.Get(instID); idRegistered {
		return true
	}

	idRegistered, _, err := core.
		getUserInstanceStateByID(instID)
	if err != nil {
		panic(err)
	}

	if idRegistered {
		core.registeredUserIDCache.Add(instID, nil)
	}

	return idRegistered
}
