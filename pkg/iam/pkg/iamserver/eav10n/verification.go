package eav10n

import (
	"time"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

type VerificationMethod int

const (
	VerificationMethodUnspecified VerificationMethod = iota
	VerificationMethodUnknown

	VerificationMethodNone
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
	}
	return VerificationMethodUnknown
}

//TODO: make this private
type verificationDBModel struct {
	IDNum                         int64              `db:"id_num"`
	Local                         string             `db:"local_part"`
	Domain                        string             `db:"domain_part"`
	Code                          string             `db:"code"`
	CodeExpiry                    *time.Time         `db:"code_expiry"`
	CreationTime                  time.Time          `db:"md_c_ts"`
	CreationUserIDNum             *iam.UserIDNum     `db:"md_c_uid"`
	CreationTerminalIDNum         *iam.TerminalIDNum `db:"md_c_tid"`
	ConfirmationAttemptsRemaining int16              `db:"confirmation_attempts_remaining"`
	ConfirmationTime              *time.Time         `db:"confirmation_ts"`
	ConfirmationUserIDNum         *iam.UserIDNum     `db:"confirmation_uid"`
	ConfirmationTerminalIDNum     *iam.TerminalIDNum `db:"confirmation_tid"`
}
