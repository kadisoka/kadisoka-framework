package iam

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"
)

type UserService interface {
	UserInstanceService
	UserProfileService

	UserKeyPhoneNumberService
	UserKeyEmailAddressService
}

var (
	ErrUserKeyPhoneNumberConflict = errors.EntMsg("user key phone number", "conflict")
)

type UserKeyPhoneNumber struct {
	UserRef     UserRefKey
	PhoneNumber PhoneNumber
}

// JSONV1 models

type UserPhoneNumberJSONV1 struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
}

type UserPhoneNumberListJSONV1 struct {
	Items []UserPhoneNumberJSONV1 `json:"items"`
}

type UserContactListsJSONV1 struct {
	Items []UserJSONV1 `json:"items"`
}
