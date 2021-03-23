package eav10n

import (
	"time"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

type VerificationMethod int

const (
	VerificationMethodUnspecified VerificationMethod = iota
	VerificationMethodNone
)

func VerificationMethodFromString(str string) VerificationMethod {
	switch str {
	case "none":
		return VerificationMethodNone
	}
	return VerificationMethodUnspecified
}

//TODO: make this private
type verificationDBModel struct {
	ID                     int64              `db:"id"`
	Local                  string             `db:"local_part"`
	Domain                 string             `db:"domain_part"`
	Code                   string             `db:"code"`
	CodeExpiry             *time.Time         `db:"code_expiry"`
	AttemptsRemaining      int16              `db:"attempts_remaining"`
	CreationTime           time.Time          `db:"c_ts"`
	CreationUserID         *iam.UserIDNum     `db:"c_uid"`
	CreationTerminalID     *iam.TerminalIDNum `db:"c_tid"`
	ConfirmationTime       *time.Time         `db:"confirmation_ts"`
	ConfirmationUserID     *iam.UserIDNum     `db:"confirmation_uid"`
	ConfirmationTerminalID *iam.TerminalIDNum `db:"confirmation_tid"`
}
