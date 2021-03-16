package iamserver

import (
	"time"

	"golang.org/x/text/language"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"
)

//TODO: some data should be taken from the context instead of
// provided in here.
type TerminalRegistrationInput struct {
	ApplicationRef iam.ApplicationRefKey
	UserRef        iam.UserRefKey

	DisplayName    string
	AcceptLanguage string //TODO: remove this

	VerificationType string
	VerificationID   int64
	VerificationTime *time.Time
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
	DisplayName            string
	UserAgentString        string
	UserPreferredLanguages []language.Tag
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
	ID       iam.TerminalID    `db:"id"`
	ClientID iam.ApplicationID `db:"application_id"`
	UserID   iam.UserID        `db:"user_id"`

	CreationTime          time.Time       `db:"c_ts"`
	CreationUserID        *iam.UserID     `db:"c_uid"`
	CreationTerminalID    *iam.TerminalID `db:"c_tid"`
	CreationOriginAddress string          `db:"c_origin_address"`
	CreationOriginEnv     string          `db:"c_origin_env"`

	DeletionTime       *time.Time      `db:"d_ts"`
	DeletionUserID     *iam.UserID     `db:"d_uid"`
	DeletionTerminalID *iam.TerminalID `db:"d_tid"`

	Secret         string `db:"secret"`
	DisplayName    string `db:"display_name"`
	AcceptLanguage string `db:"accept_language"`

	VerificationType string     `db:"verification_type"`
	VerificationID   int64      `db:"verification_id"`
	VerificationTime *time.Time `db:"verification_ts"`
}
