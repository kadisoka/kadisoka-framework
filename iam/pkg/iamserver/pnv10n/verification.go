package pnv10n

import (
	"time"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

type VerificationMethod int

const (
	VerificationMethodUnspecified VerificationMethod = iota
	VerificationMethodUnknown

	VerificationMethodNone
	VerificationMethodSMS
)

func (method VerificationMethod) IsValid() bool {
	return method != VerificationMethodUnspecified &&
		method != VerificationMethodUnknown
}

func VerificationMethodFromString(str string) VerificationMethod {
	switch str {
	case "":
		return VerificationMethodUnspecified
	case "none":
		return VerificationMethodNone
	case "sms":
		return VerificationMethodSMS
	}
	return VerificationMethodUnknown
}

//TODO: make this private
type verificationDBModel struct {
	IDNum                     int64              `db:"id"`
	CountryCode               int32              `db:"country_code"`
	NationalNumber            int64              `db:"national_number"`
	Code                      string             `db:"code"`
	CodeExpiry                *time.Time         `db:"code_expiry"`
	AttemptsRemaining         int16              `db:"attempts_remaining"`
	CreationTime              time.Time          `db:"c_ts"`
	CreationUserIDNum         *iam.UserIDNum     `db:"c_uid"`
	CreationTerminalIDNum     *iam.TerminalIDNum `db:"c_tid"`
	ConfirmationTime          *time.Time         `db:"confirmation_ts"`
	ConfirmationUserIDNum     *iam.UserIDNum     `db:"confirmation_uid"`
	ConfirmationTerminalIDNum *iam.TerminalIDNum `db:"confirmation_tid"`
}
