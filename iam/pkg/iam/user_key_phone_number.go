package iam

import "github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"

// Key phone number is a phone number which can be used to sign in.

type UserKeyPhoneNumberService interface {
	GetUserKeyPhoneNumber(
		callCtx CallContext,
		userRef UserRefKey,
	) (*telephony.PhoneNumber, error)
}
