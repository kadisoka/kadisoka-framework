package iam

import (
	"crypto/rand"
	"encoding/binary"
	"strings"

	azcore "github.com/alloyzeus/go-azfl/azcore"
	azid "github.com/alloyzeus/go-azfl/azid"
	errors "github.com/alloyzeus/go-azfl/errors"
	email "github.com/kadisoka/kadisoka-framework/volib/pkg/email"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the azfl package it is being compiled against.
// A compilation error at this line likely means your copy of the
// azfl package needs to be updated.
var _ = azcore.AZCorePackageIsVersion1

// Reference imports to suppress errors if they are not otherwise used.
var _ = azid.BinDataTypeUnspecified
var _ = errors.ErrUnimplemented
var _ = binary.MaxVarintLen16
var _ = rand.Reader
var _ = strings.Compare

// Adjunct-prime UserKeyEmailAddress of User.
//
// A key email address is an email address that can be used as the identifier
// for logging in by the specified user.

//region Service

// UserKeyEmailAddressService provides a contract
// for methods related to adjunct UserKeyEmailAddress.
type UserKeyEmailAddressService interface {
	// AZxAdjuncPrimeService

	GetUserKeyEmailAddress(
		inputCtx CallInputContext,
		userID UserID,
	) (*email.Address, error)
}

//endregion
