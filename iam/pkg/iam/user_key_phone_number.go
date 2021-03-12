package iam

// Key phone number is a phone number which can be used to sign in.

type UserKeyPhoneNumberService interface {
	GetUserKeyPhoneNumber(
		callCtx CallContext,
		userRef UserRefKey,
	) (*PhoneNumber, error)
}
