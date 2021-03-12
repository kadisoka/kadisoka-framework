package iam

// Key email address is an email address which can be used to sign in.

type UserKeyEmailAddressService interface {
	GetUserKeyEmailAddress(
		callCtx CallContext,
		userRef UserRefKey,
	) (*EmailAddress, error)
}
