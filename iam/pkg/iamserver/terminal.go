package iamserver

import (
	"time"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
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

// terminalDBModel represents a row from terminal table.
type terminalDBModel struct {
	ID       iam.TerminalID    `db:"id"`
	ClientID iam.ApplicationID `db:"application_id"`
	UserID   iam.UserID        `db:"user_id"`
	Secret   string            `db:"secret"`

	DisplayName    string `db:"display_name"`
	AcceptLanguage string `db:"accept_language"`

	CreationTime       time.Time       `db:"c_ts"`
	CreationUserID     *iam.UserID     `db:"c_uid"`
	CreationTerminalID *iam.TerminalID `db:"c_tid"`
	CreationIPAddress  string          `db:"c_origin_address"`
	CreationUserAgent  string          `db:"c_origin_env"`

	VerificationType string     `db:"verification_type"`
	VerificationID   int64      `db:"verification_id"`
	VerificationTime *time.Time `db:"verification_time"`
}
