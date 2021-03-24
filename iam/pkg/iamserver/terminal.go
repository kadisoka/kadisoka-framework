package iamserver

import (
	"time"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
)

//TODO: some data should be taken from the context instead of
// provided in here.
type TerminalRegistrationInput struct {
	Context        iam.CallContext
	ApplicationRef iam.ApplicationRefKey //TODO: get from context
	Data           TerminalRegistrationInputData
}

type TerminalRegistrationInputData struct {
	UserRef iam.UserRefKey

	DisplayName string

	VerificationType string
	VerificationID   int64
	VerificationTime *time.Time
}

type TerminalRegistrationOutput struct {
	Context iam.OpOutputContext
	Data    TerminalRegistrationOutputData
}

type TerminalRegistrationOutputData struct {
	TerminalRef    iam.TerminalRefKey
	TerminalSecret string
}

//TODO: use generics when it's available
type TerminalAuthorizationByEmailAddressStartInput struct {
	Context        iam.CallContext
	ApplicationRef iam.ApplicationRefKey //TODO: should be from Context.Authorization
	Data           TerminalAuthorizationByEmailAddressStartInputData
}

type TerminalAuthorizationByEmailAddressStartInputData struct {
	EmailAddress        iam.EmailAddress
	VerificationMethods []eav10n.VerificationMethod

	TerminalAuthorizationStartInputBaseData
}

//TODO: use generics when it's available
type TerminalAuthorizationByPhoneNumberStartInput struct {
	Context        iam.CallContext
	ApplicationRef iam.ApplicationRefKey //TODO: should be from Context.Authorization
	Data           TerminalAuthorizationByPhoneNumberStartInputData
}

type TerminalAuthorizationByPhoneNumberStartInputData struct {
	PhoneNumber         iam.PhoneNumber
	VerificationMethods []pnv10n.VerificationMethod

	TerminalAuthorizationStartInputBaseData
}

type TerminalAuthorizationStartInputBaseData struct {
	DisplayName string
}

type TerminalAuthorizationStartOutput struct {
	Context iam.OpOutputContext
	Data    TerminalAuthorizationStartOutputData
}

type TerminalAuthorizationStartOutputData struct {
	TerminalRef                iam.TerminalRefKey
	VerificationID             int64
	VerificationCodeExpiryTime *time.Time
}

// terminalDBRawModel represents a row from terminal table.
type terminalDBRawModel struct {
	IDNum            iam.TerminalIDNum    `db:"id"`
	ApplicationIDNum iam.ApplicationIDNum `db:"application_id"`
	UserIDNum        iam.UserIDNum        `db:"user_id"`

	CreationTime          time.Time          `db:"c_ts"`
	CreationUserIDNum     *iam.UserIDNum     `db:"c_uid"`
	CreationTerminalIDNum *iam.TerminalIDNum `db:"c_tid"`
	CreationOriginAddress string             `db:"c_origin_address"`
	CreationOriginEnv     string             `db:"c_origin_env"`

	DeletionTime          *time.Time         `db:"d_ts"`
	DeletionUserIDNum     *iam.UserIDNum     `db:"d_uid"`
	DeletionTerminalIDNum *iam.TerminalIDNum `db:"d_tid"`

	Secret         string `db:"secret"`
	DisplayName    string `db:"display_name"`
	AcceptLanguage string `db:"accept_language"`

	VerificationType string     `db:"verification_type"`
	VerificationID   int64      `db:"verification_id"`
	VerificationTime *time.Time `db:"verification_ts"`
}
