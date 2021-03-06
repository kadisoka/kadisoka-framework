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
	ID                     int64           `db:"id"`
	Local                  string          `db:"local_part"`
	Domain                 string          `db:"domain_part"`
	Code                   string          `db:"code"`
	CodeExpiry             *time.Time      `db:"code_expiry"`
	AttemptsRemaining      int16           `db:"attempts_remaining"`
	CreationTime           time.Time       `db:"c_ts"`
	CreationUserID         *iam.UserID     `db:"c_uid"`
	CreationTerminalID     *iam.TerminalID `db:"c_tid"`
	ConfirmationTime       *time.Time      `db:"confirmation_time"`
	ConfirmationUserID     *iam.UserID     `db:"confirmation_user_id"`
	ConfirmationTerminalID *iam.TerminalID `db:"confirmation_terminal_id"`
}
