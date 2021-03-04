package iamserver

import (
	"time"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

//TODO: some data should be taken from the context instead of
// provided in here.
type TerminalRegistrationInput struct {
	ClientID iam.ClientID
	UserRef  iam.UserRefKey

	DisplayName    string
	AcceptLanguage string //TODO: remove this

	VerificationType string
	VerificationID   int64
	VerificationTime *time.Time
}

// terminalDBModel represents a row from 'terminals' table.
type terminalDBModel struct {
	ID       iam.TerminalID `db:"id"`
	ClientID iam.ClientID   `db:"client_id"`
	UserID   iam.UserID     `db:"user_id"`
	Secret   string         `db:"secret"`

	DisplayName    string `db:"display_name"`
	AcceptLanguage string `db:"accept_language"`
	PlatformType   string `db:"platform_type"` //TODO(exa): remove this. get from client info.

	CreationTime       time.Time       `db:"creation_time"`
	CreationUserID     *iam.UserID     `db:"creation_user_id"`
	CreationTerminalID *iam.TerminalID `db:"creation_terminal_id"`
	CreationIPAddress  string          `db:"creation_ip_address"`
	CreationUserAgent  string          `db:"creation_user_agent"`

	VerificationType string     `db:"verification_type"`
	VerificationID   int64      `db:"verification_id"`
	VerificationTime *time.Time `db:"verification_time"`
}
