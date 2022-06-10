package iam

import "github.com/kadisoka/kadisoka-framework/volib/pkg/email"

// Key email address is an email address which can be used to sign in.

type UserKeyEmailAddressService interface {
	GetUserKeyEmailAddress(
		inputCtx CallInputContext,
		userID UserID,
	) (*email.Address, error)
}
