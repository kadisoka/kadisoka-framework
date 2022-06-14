package iamserver

import (
	"time"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver/pnv10n"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/email"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/telephony"
)

type TerminalRegistrationInputData struct {
	// ApplicationID is the ID of the authenticated application that
	// the new terminal will be associated to.
	ApplicationID iam.ApplicationID
	UserID        iam.UserID

	DisplayName string

	VerificationType string
	VerificationID   int64
	VerificationTime *time.Time
}

type TerminalRegistrationOutputData struct {
	TerminalID     iam.TerminalID
	TerminalSecret string
}

type TerminalAuthorizationByEmailAddressStartInputData struct {
	EmailAddress        email.Address
	VerificationMethods []eav10n.VerificationMethod

	TerminalAuthorizationStartInputBaseData
}

type TerminalAuthorizationByPhoneNumberStartInputData struct {
	PhoneNumber         telephony.PhoneNumber
	VerificationMethods []pnv10n.VerificationMethod

	TerminalAuthorizationStartInputBaseData
}

type TerminalAuthorizationStartInputBaseData struct {
	// ApplicationID is the ID of the authenticated application that
	// the new terminal will be associated to.
	ApplicationID iam.ApplicationID

	DisplayName string
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
