package iam

import "github.com/kadisoka/kadisoka-framework/volib/pkg/email"

// Key email address is an email address which can be used to sign in.

type UserKeyEmailAddressService interface {
	GetUserKeyEmailAddress(
		callCtx CallContext,
		userRef UserRefKey,
	) (*email.Address, error)
}
