package iam

import (
	"github.com/alloyzeus/go-azcore/azcore/errors"
)

type UserService interface {
	UserAccountService
	UserProfileService

	GetUserIdentifierPhoneNumber(
		callCtx CallContext,
		userID UserID,
	) (*PhoneNumber, error)

	GetUserIdentifierEmailAddress(
		callCtx CallContext,
		userID UserID,
	) (*EmailAddress, error)
}

//TODO: this does not belong to C2S service, but only in S2S service
type UserTerminalService interface {
	ListUserTerminalIDFirebaseInstanceTokens(
		ownerUserID UserID,
	) ([]TerminalIDFirebaseInstanceToken, error)
	DeleteUserTerminalFCMRegistrationToken(
		authCtx *Authorization,
		userID UserID, terminalID TerminalID, token string,
	) error
}

var (
	ErrUserIdentifierPhoneNumberConflict = errors.EntMsg("user identifier phone number", "conflict")
)

type UserIdentifierPhoneNumber struct {
	UserID      UserID
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

type UserEmailAddressPutRequestJSONV1 struct {
	IsPrimary bool `json:"is_primary" db:"is_primary"`
}

type UserContactListsJSONV1 struct {
	Items []UserJSONV1 `json:"items"`
}
