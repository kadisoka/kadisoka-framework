package iamserver

import "github.com/kadisoka/kadisoka-framework/iam/pkg/iam"

// Interface conformance assertion.
var _ iam.UserIDService = &Core{}

// IsUserIDRegistered is used to determine that a user ID has been registered.
// It's not checking if the account is active or not.
//
// This function is generally cheap if the user ID has been registered.
func (core *Core) IsUserIDRegistered(id iam.UserID) bool {
	// Look up for an user ID in the cache.
	if _, idRegistered := core.registeredUserIDCache.Get(id); idRegistered {
		return true
	}

	idRegistered, _, err := core.
		getUserAccountState(id)
	if err != nil {
		panic(err)
	}

	if idRegistered {
		core.registeredUserIDCache.Add(id, nil)
	}

	return idRegistered
}
