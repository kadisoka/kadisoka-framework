package iamserver

import (
	"time"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/email"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

//TODO: some data should be taken from the context instead of
// provided in here.
type TerminalRegistrationInput struct {
	Context       iam.CallInputContext
	ApplicationID iam.ApplicationID //TODO: get from context
	Data          TerminalRegistrationInputData
}

type TerminalRegistrationInputData struct {
	UserID iam.UserID

	DisplayName string

	VerificationType string
	VerificationID   int64
	VerificationTime *time.Time
}

type TerminalRegistrationOutput struct {
	Context iam.CallOutputContext
	Data    TerminalRegistrationOutputData
}

type TerminalRegistrationOutputData struct {
	TerminalID     iam.TerminalID
	TerminalSecret string
}

//TODO: use generics when it's available
type TerminalAuthorizationByEmailAddressStartInput struct {
	Context       iam.CallInputContext
	ApplicationID iam.ApplicationID //TODO: should be from Context.Authorization
	Data          TerminalAuthorizationByEmailAddressStartInputData
}

type TerminalAuthorizationByEmailAddressStartInputData struct {
	EmailAddress        email.Address
	VerificationMethods []eav10n.VerificationMethod

	TerminalAuthorizationStartInputBaseData
}

//TODO: use generics when it's available
type TerminalAuthorizationByPhoneNumberStartInput struct {
	Context       iam.CallInputContext
	ApplicationID iam.ApplicationID //TODO: should be from Context.Authorization
	Data          TerminalAuthorizationByPhoneNumberStartInputData
}

type TerminalAuthorizationByPhoneNumberStartInputData struct {
	PhoneNumber         telephony.PhoneNumber
	VerificationMethods []pnv10n.VerificationMethod

	TerminalAuthorizationStartInputBaseData
}

type TerminalAuthorizationStartInputBaseData struct {
	DisplayName string
}

type TerminalAuthorizationStartOutput struct {
	Context iam.CallOutputContext
	Data    TerminalAuthorizationStartOutputData
}

type TerminalAuthorizationStartOutputData struct {
	TerminalID                 iam.TerminalID
	VerificationID             int64
	VerificationCodeExpiryTime *time.Time
}

// terminalDBRawModel represents a row from terminal table.
type terminalDBRawModel struct {
	IDNum            iam.TerminalIDNum    `db:"id_num"`
	ApplicationIDNum iam.ApplicationIDNum `db:"application_id"`
	UserIDNum        iam.UserIDNum        `db:"user_id"`

	CreationTime          time.Time          `db:"_mc_ts"`
	CreationUserIDNum     *iam.UserIDNum     `db:"_mc_uid"`
	CreationTerminalIDNum *iam.TerminalIDNum `db:"_mc_tid"`
	CreationOriginAddress string             `db:"_mc_origin_address"`
	CreationOriginEnv     string             `db:"_mc_origin_env"`

	DeletionTime          *time.Time         `db:"_md_ts"`
	DeletionUserIDNum     *iam.UserIDNum     `db:"_md_uid"`
	DeletionTerminalIDNum *iam.TerminalIDNum `db:"_md_tid"`

	Secret         string `db:"secret"`
	DisplayName    string `db:"display_name"`
	AcceptLanguage string `db:"accept_language"`

	VerificationType string     `db:"verification_type"`
	VerificationID   int64      `db:"verification_id"`
	VerificationTime *time.Time `db:"verification_ts"`
}
