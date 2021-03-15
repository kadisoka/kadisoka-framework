package iam

import (
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	"golang.org/x/text/language"
)

type TerminalService interface {
	GetTerminalInfo(
		callCtx CallContext,
		terminalID TerminalID,
	) (*TerminalInfo, error)
}

const (
	TerminalVerificationResourceTypePhoneNumber  = "phone-number"
	TerminalVerificationResourceTypeEmailAddress = "email-address"

	TerminalVerificationResourceTypeOAuthAuthorizationCode = "oauth2-authorization-code"
	TerminalVerificationResourceTypeOAuthImplicit          = "oauth2-implicit"
	TerminalVerificationResourceTypeOAuthClientCredentials = "oauth2-client-credentials"
)

var (
	ErrTerminalVerificationCodeMismatch = errors.EntMsg("terminal verification code", "mismatch")
	ErrTerminalVerificationCodeExpired  = errors.EntMsg("terminal verification code", "expired")

	ErrTerminalVerificationResourceConflict = errors.EntMsg("terminal verification resource", "conflict")

	ErrTerminalVerificationResourceNameInvalid = errors.Ent("terminal verification resource name", nil)
)

type TerminalInfo struct {
	DisplayName    string
	AcceptLanguage []language.Tag
}

//TODO: this does not belong to C2S service, but only in S2S service
type TerminalFCMRegistrationTokenService interface {
	ListTerminalFCMRegistrationTokensByUser(
		ownerUserRef UserRefKey,
	) (tokens map[TerminalID]string, err error)
	DisposeTerminalFCMRegistrationToken(
		callCtx CallContext,
		terminalID TerminalID,
		token string,
	) error
}

// JSONV1 models

type TerminalRegisterPostRequestJSONV1 struct {
	VerificationResourceName string   `json:"verification_resource_name"`
	VerificationMethods      []string `json:"verification_methods"`
	DisplayName              string   `json:"display_name"`
}

func (TerminalRegisterPostRequestJSONV1) SwaggerDoc() map[string]string {
	return map[string]string{
		"verification_resource_name": "A phone number complete with country code or an email address.",
		"verification_methods": "The preferred verification methods. " +
			"The values are resource-type-specific. For phone-number, it defaults to SMS.",
		"display_name": "For the user to make it easy to identify. " +
			"The recommended value is the user's device name.",
	}
}

// provide user id? indicator for a new user?
type TerminalRegisterPostResponseJSONV1 struct {
	TerminalID     string     `json:"terminal_id"`
	TerminalSecret string     `json:"terminal_secret,omitempty"`
	CodeExpiry     *time.Time `json:"code_expiry,omitempty"`
}

func (TerminalRegisterPostResponseJSONV1) SwaggerDoc() map[string]string {
	return map[string]string{
		"terminal_id": "The ID for the terminal.",
		"terminal_secret": "Contains terminal's secret for certain " +
			"verification resource types",
		"code_expiry": "The time when the verification code will " +
			"be expired.",
	}
}
